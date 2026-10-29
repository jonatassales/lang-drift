// Package manifest reads a component's manifests from disk and hands them to
// the rules that classify them.
package manifest

import (
	"errors"
	"io/fs"
	"path/filepath"

	domain "github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// manifestLimit bounds one manifest; a manifest larger than this is a lie.
const manifestLimit = 1 << 20

// Inspector reads manifests from the project.
type Inspector struct{}

// Inspect reads the known manifests under a component path. A path without
// any is an empty inspection, not an error.
func (Inspector) Inspect(owner project.Project, componentPath string) (domain.Inspection, error) {
	directory, err := files.Existing(owner.Root, filepath.FromSlash(componentPath))
	if err != nil {
		return domain.Inspection{}, err
	}

	manifests := map[string][]byte{}
	for _, name := range domain.Names {
		path, err := files.Existing(directory, name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return domain.Inspection{}, err
		}

		content, err := files.ReadBounded(path, manifestLimit)
		if err != nil {
			return domain.Inspection{}, err
		}

		manifests[name] = content
	}

	return domain.Classify(manifests), nil
}
