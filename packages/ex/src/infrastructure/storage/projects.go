package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// projectMetadataVersion is the only project metadata format EX reads.
const projectMetadataVersion = 1

// Directory permissions: metadata is shared with the consumer; state is private.
const (
	metadataDirectoryMode = 0o755
	stateDirectoryMode    = 0o700
	stateIgnoreMode       = 0o600
)

type projectMetadata struct {
	Version int    `json:"version"`
	ID      string `json:"id"`
}

// Projects keeps project identity in .exphere/config.json.
type Projects struct{}

// Open reads the project at the selected directory, or the nearest one above it.
func (Projects) Open(selection application.ProjectSelection) (project.Project, error) {
	start, err := files.CanonicalDirectory(selection.Directory)
	if err != nil {
		return project.Project{}, err
	}

	if selection.Explicit {
		return openAt(start)
	}

	return discover(start)
}

// Create records a new project identity. It never replaces existing metadata
// and never adopts existing state.
func (Projects) Create(directory string, id project.ID) (project.Project, error) {
	root, err := files.CanonicalDirectory(directory)
	if err != nil {
		return project.Project{}, err
	}

	layout := layoutAt(root)

	if _, err := files.Within(root, filepath.Join(".exphere", "state")); err != nil {
		return project.Project{}, err
	}

	if err := os.MkdirAll(layout.metadataDirectory(), metadataDirectoryMode); err != nil {
		return project.Project{}, err
	}

	if err := os.MkdirAll(layout.state(), stateDirectoryMode); err != nil {
		return project.Project{}, err
	}

	unlock, err := files.Lock(layout.initializationLock())
	if err != nil {
		return project.Project{}, err
	}
	defer unlock()

	if _, err := os.Lstat(layout.metadata()); err == nil {
		return project.Project{}, application.AlreadyExists("project metadata already exists at %s", root)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return project.Project{}, err
	}

	if err := ensureStateIgnored(layout.stateIgnore()); err != nil {
		return project.Project{}, err
	}

	metadata := projectMetadata{Version: projectMetadataVersion, ID: id.String()}
	if err := files.WriteJSON(layout.metadata(), metadata); err != nil {
		return project.Project{}, err
	}

	return project.Project{ID: id, Root: root}, nil
}

// discover walks upward to the nearest project metadata. Invalid metadata
// stops the walk rather than silently selecting an ancestor.
func discover(start string) (project.Project, error) {
	for root := start; ; root = filepath.Dir(root) {
		_, err := os.Lstat(layoutAt(root).metadata())

		switch {
		case err == nil:
			return openAt(root)

		case !errors.Is(err, fs.ErrNotExist):
			return project.Project{}, err

		case root == filepath.Dir(root):
			return project.Project{}, application.NotFound(
				"no Exphere project at or above %s; initialize one at the project root",
				start,
			)
		}
	}
}

// openAt reads the project at exactly root. Unsupported metadata fails
// explicitly rather than being migrated or treated as an empty project.
func openAt(root string) (project.Project, error) {
	var metadata projectMetadata

	err := files.ReadJSON(layoutAt(root).metadata(), &metadata)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return project.Project{}, application.NotFound("no Exphere project at %s", root)

	case err != nil:
		return project.Project{}, fmt.Errorf("open project: %w", err)

	case metadata.Version != projectMetadataVersion:
		return project.Project{}, fmt.Errorf("unsupported project metadata version %d", metadata.Version)
	}

	id, err := project.ParseID(metadata.ID)
	if err != nil {
		return project.Project{}, err
	}

	if _, err := files.Within(root, filepath.Join(".exphere", "state")); err != nil {
		return project.Project{}, err
	}

	return project.Project{ID: id, Root: root}, nil
}

// ensureStateIgnored keeps private state out of version control from inside
// the state directory, so the consumer's own ignore rules are never edited.
func ensureStateIgnored(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, stateIgnoreMode)
	if err == nil {
		_, writeErr := file.WriteString("*\n")

		return errors.Join(writeErr, file.Close())
	}

	if !errors.Is(err, fs.ErrExist) {
		return err
	}

	contents, err := files.ReadBounded(path, 64)
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(contents)) != "*" {
		return errors.New("the existing state ignore file must ignore every private file")
	}

	return nil
}
