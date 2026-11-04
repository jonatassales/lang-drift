package scaffold

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
)

// ciGatesTemplate is the workflow as ex ships it, with {{version}} where the
// generating version goes.
const ciGatesTemplate = "assets/github/langdrift-gates.yml"

//go:embed assets/github/langdrift-gates.yml
var ciGatesAssets embed.FS

// ciGatesWorkflow is where the CI Gates workflow lands in the project: the
// one file ex init may write outside .langdrift/, after an explicit yes.
const ciGatesWorkflow = ".github/workflows/langdrift-gates.yml"

// CIGates writes the CI Gates workflow into the project. It is additive: an
// existing file is kept byte for byte, never replaced.
type CIGates struct{}

// Exists reports whether the project already has the workflow file.
func (gates CIGates) Exists(owner project.Project) (bool, error) {
	target, err := gates.target(owner)
	if err != nil {
		return false, err
	}

	_, err = os.Lstat(target)

	switch {
	case err == nil:
		return true, nil

	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	}

	return false, err
}

// Write writes the workflow pinned to version, unless the file appeared in
// the meantime; it reports whether it wrote it.
func (gates CIGates) Write(owner project.Project, version string) (bool, error) {
	target, err := gates.target(owner)
	if err != nil {
		return false, err
	}

	template, err := ciGatesAssets.ReadFile(ciGatesTemplate)
	if err != nil {
		return false, err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return false, err
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	content := strings.ReplaceAll(string(template), "{{version}}", version)

	_, writeErr := file.WriteString(content)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return false, fmt.Errorf("write %s: %w", ciGatesWorkflow, err)
	}

	return true, nil
}

// target is the workflow's path inside the project.
func (CIGates) target(owner project.Project) (string, error) {
	return files.Within(owner.Root, filepath.FromSlash(ciGatesWorkflow))
}
