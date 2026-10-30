package scaffold

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

//go:embed assets/langdrift/microruntimes/AGENTS.md
var langdriftAssets embed.FS

// The execution rules every run follows, as the engine ships them.
const executionRules = "microruntimes/AGENTS.md"

// LangDrift adds the Microruntimes execution rules to a project's .langdrift
// folder. It is additive: an existing file is kept byte for byte, never
// replaced, and nothing else is written. LangDrift prescribes no workflow
// and no artifact inside .langdrift, so no template, skill or policy text
// lands here; the CI Gates workflow is CIGates's, written only on a yes.
type LangDrift struct{}

// Prepare writes the execution rules when the project has none.
func (LangDrift) Prepare(owner project.Project) error {
	target, err := files.Within(owner.Root, filepath.Join(".langdrift", executionRules))
	if err != nil {
		return err
	}

	_, err = os.Lstat(target)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	content, err := langdriftAssets.ReadFile("assets/langdrift/" + executionRules)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}

	_, writeErr := file.Write(content)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("write %s: %w", executionRules, err)
	}

	return nil
}
