package scaffold

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

//go:embed all:assets/connected
var connectedAssets embed.FS

// Connected supplies the standalone project agent's reference workflow. It is
// additive: existing customizations are kept byte-for-byte, never merged away.
type Connected struct{}

func (Connected) Prepare(owner project.Project) error {
	return fs.WalkDir(connectedAssets, "assets/connected", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel("assets/connected", path)
		if err != nil {
			return err
		}
		target, err := files.Within(owner.Root, filepath.Join(".exphere", relative))
		if err != nil {
			return err
		}
		if _, err := os.Lstat(target); err == nil {
			return nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		content, err := connectedAssets.ReadFile(path)
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
			return fmt.Errorf("bootstrap %s: %w", relative, err)
		}
		return nil
	})
}
