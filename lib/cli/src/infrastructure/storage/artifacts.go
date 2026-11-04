package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonatassales/tao/lib/cli/src/application/inputs"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
)

// Artifacts keeps the folder of the team's tickets and specs in the project
// identity, .langdrift/config.json.
type Artifacts struct{}

// Folder is the recorded artifacts folder, or empty when none is.
func (Artifacts) Folder(owner project.Project) (string, error) {
	var metadata projectMetadata

	if err := files.ReadJSON(layoutAt(owner.Root).metadata(), &metadata); err != nil {
		return "", fmt.Errorf("read the project identity: %w", err)
	}

	return metadata.Artifacts, nil
}

// Check returns relative as it would be stored, cleaned and slash-separated,
// once it names a folder inside the project; it records nothing.
func (Artifacts) Check(owner project.Project, relative string) (string, error) {
	cleaned := filepath.Clean(relative)

	folder, err := files.Existing(owner.Root, cleaned)
	if err != nil {
		return "", fmt.Errorf("%w: the artifacts folder %q must be a folder inside the project: %v", inputs.ErrInvalid, relative, err)
	}

	info, err := os.Stat(folder)
	if err != nil {
		return "", fmt.Errorf("read the artifacts folder: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%w: the artifacts folder %q is not a folder", inputs.ErrInvalid, relative)
	}

	return filepath.ToSlash(cleaned), nil
}

// Record stores relative, once Check accepts it, and returns it as stored.
func (artifacts Artifacts) Record(owner project.Project, relative string) (string, error) {
	stored, err := artifacts.Check(owner, relative)
	if err != nil {
		return "", err
	}

	layout := layoutAt(owner.Root)

	unlock, err := files.Lock(layout.initializationLock())
	if err != nil {
		return "", err
	}
	defer unlock()

	var metadata projectMetadata

	if err := files.ReadJSON(layout.metadata(), &metadata); err != nil {
		return "", fmt.Errorf("read the project identity: %w", err)
	}

	metadata.Artifacts = stored

	encoded, err := json.MarshalIndent(metadata, "", "\t")
	if err != nil {
		return "", fmt.Errorf("encode the project identity: %w", err)
	}

	if err := files.WriteDocument(layout.metadata(), append(encoded, '\n'), metadataFileMode); err != nil {
		return "", fmt.Errorf("write the project identity: %w", err)
	}

	return metadata.Artifacts, nil
}
