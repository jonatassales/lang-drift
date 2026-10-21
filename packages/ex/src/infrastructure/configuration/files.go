package configuration

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

const (
	descriptorFileName = "schema.json"

	// maxDescriptorBytes bounds the descriptor; maxNativeFileBytes bounds each
	// manifest and Terraform source file.
	maxDescriptorBytes = 4 << 20
	maxNativeFileBytes = 1 << 20
)

// manifestExtensions are the native manifest formats EX reads.
var manifestExtensions = map[string]bool{".json": true, ".yaml": true, ".yml": true}

// Files reads the configuration consumers author under .exphere/runtimes.
// It checks that references exist and parses native files; the policy those
// files must satisfy belongs to the domain.
type Files struct{}

// Load reads a runtime's descriptor exactly as authored.
func (Files) Load(owner project.Project, runtime microruntime.Identifier) (microruntime.Descriptor, error) {
	if err := rejectLegacyRuntime(owner.Root, runtime); err != nil {
		return microruntime.Descriptor{}, err
	}

	directory, err := runtimeDirectory(owner, runtime)
	if err != nil {
		return microruntime.Descriptor{}, err
	}

	path, err := files.Existing(directory, descriptorFileName)
	if err != nil {
		return microruntime.Descriptor{}, err
	}

	data, err := files.ReadBounded(path, maxDescriptorBytes)
	if err != nil {
		return microruntime.Descriptor{}, err
	}

	return decodeDescriptor(data)
}

// Snapshot reads every native file a valid descriptor references and
// fingerprints them together with the descriptor.
func (Files) Snapshot(owner project.Project, descriptor microruntime.Descriptor) (application.Snapshot, error) {
	directory, err := runtimeDirectory(owner, descriptor.ID)
	if err != nil {
		return application.Snapshot{}, err
	}

	if err := requireComponents(owner.Root, descriptor); err != nil {
		return application.Snapshot{}, err
	}

	if err := requireHooks(directory, descriptor); err != nil {
		return application.Snapshot{}, err
	}

	manifests, manifestDigest, previousDigest, err := readManifests(directory, descriptor)
	if err != nil {
		return application.Snapshot{}, err
	}

	preparation, err := files.Existing(directory, descriptor.TerraformDirectory.String())
	if err != nil {
		return application.Snapshot{}, fmt.Errorf("terraform: %w", err)
	}

	terraformDigests, err := validateTerraformRoot(preparation)
	if err != nil {
		return application.Snapshot{}, fmt.Errorf("terraform: %w", err)
	}

	fingerprint, err := fingerprintOf(manifestDigest, terraformDigests)
	if err != nil {
		return application.Snapshot{}, err
	}
	previousFingerprint, err := fingerprintOf(previousDigest, terraformDigests)
	if err != nil {
		return application.Snapshot{}, err
	}

	return application.Snapshot{
		Manifests:            manifests,
		PreparationDirectory: preparation,
		Hash:                 fingerprint,
		CompatibleHashes:     []string{previousFingerprint},
	}, nil
}

func runtimeDirectory(owner project.Project, runtime microruntime.Identifier) (string, error) {
	directory, err := files.Existing(owner.Root, microruntime.ConfigurationDirectory(runtime).String())

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "", application.NotFound("runtime %s is not initialized", runtime)
	case err != nil:
		return "", err
	}

	return directory, nil
}

// requireComponents checks that every unit's source directory exists.
func requireComponents(root string, descriptor microruntime.Descriptor) error {
	for _, name := range descriptor.UnitNames() {
		if err := requireDirectory(root, descriptor.Units[name].Component); err != nil {
			return fmt.Errorf("unit %s component: %w", name, err)
		}
	}

	return nil
}

func requireDirectory(root string, relative microruntime.RelativePath) error {
	path, err := files.Existing(root, relative.String())
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s must be a directory", relative)
	}

	return nil
}

// requireHooks checks that every declared hook is a regular file.
func requireHooks(directory string, descriptor microruntime.Descriptor) error {
	for name, script := range descriptor.Hooks {
		path, err := files.Existing(directory, script.String())
		if err != nil {
			return fmt.Errorf("hook %s: %w", name, err)
		}

		if info, err := os.Lstat(path); err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("hook %s must be a regular file", name)
		}
	}

	return nil
}

// readManifests reads every unit's native objects in a stable order and
// digests the descriptor together with every manifest file read.
func readManifests(directory string, descriptor microruntime.Descriptor) ([]microruntime.Manifest, string, string, error) {
	digest := sha256.New()
	previousDigest := sha256.New()

	encoded, err := encodeDescriptor(descriptor)
	if err != nil {
		return nil, "", "", err
	}

	// Execution version 1 originally hashed compact JSON. The reviewed
	// standalone refactor used the indented scaffold bytes instead. Recompute
	// both from the same current files; neither permits changed configuration.
	var compact bytes.Buffer
	if err := json.Compact(&compact, encoded); err != nil {
		return nil, "", "", err
	}
	digest.Write(compact.Bytes())
	previousDigest.Write(encoded)

	var manifests []microruntime.Manifest

	for _, unit := range descriptor.UnitNames() {
		root, err := files.Existing(directory, descriptor.Units[unit].Manifests.String())
		if err != nil {
			return nil, "", "", fmt.Errorf("unit %s manifests: %w", unit, err)
		}

		unitManifests, err := readUnitManifests(root, unit, io.MultiWriter(digest, previousDigest))
		if err != nil {
			return nil, "", "", fmt.Errorf("unit %s: %w", unit, err)
		}

		manifests = append(manifests, unitManifests...)
	}

	return manifests, hex.EncodeToString(digest.Sum(nil)), hex.EncodeToString(previousDigest.Sum(nil)), nil
}

// readUnitManifests reads a manifest file, or every manifest file below a
// directory in lexical order. Symlinks are refused.
func readUnitManifests(root string, unit microruntime.Identifier, digest io.Writer) ([]microruntime.Manifest, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() && (!info.Mode().IsRegular() || !manifestExtensions[filepath.Ext(root)]) {
		return nil, errors.New("manifests must be a directory or a JSON or YAML file")
	}

	var manifests []microruntime.Manifest

	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		switch {
		case walkErr != nil:
			return walkErr
		case entry.Type()&fs.ModeSymlink != 0:
			return fmt.Errorf("manifest symlink refused: %s", path)
		case entry.IsDir() || !manifestExtensions[filepath.Ext(path)]:
			return nil
		}

		data, err := files.ReadBounded(path, maxNativeFileBytes)
		if err != nil {
			return err
		}

		digest.Write(data)

		objects, err := decodeObjects(data)
		if err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}

		for _, object := range objects {
			manifests = append(manifests, microruntime.Manifest{Unit: unit, Object: object})
		}

		return nil
	})

	return manifests, walkErr
}

// decodeObjects reads every non-empty document of a YAML or JSON stream.
func decodeObjects(data []byte) ([]microruntime.Object, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))

	var objects []microruntime.Object

	for {
		var object map[string]any

		err := decoder.Decode(&object)

		switch {
		case errors.Is(err, io.EOF):
			return objects, nil
		case err != nil:
			return nil, err
		case len(object) > 0:
			objects = append(objects, object)
		}
	}
}

// fingerprintOf combines the manifest digest with each Terraform file's
// digest, so any change to either invalidates cleanup.
func fingerprintOf(manifestDigest string, terraformDigests map[string]string) (string, error) {
	parts := map[string]string{"descriptor-and-manifests": manifestDigest}
	for name, digest := range terraformDigests {
		parts[name] = digest
	}

	serialized, err := json.Marshal(parts)
	if err != nil {
		return "", err
	}

	return digestOf(serialized), nil
}

func digestOf(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}
