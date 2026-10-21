package configuration

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/exphere/ex/src/domain/microruntime"
)

// rejectLegacyRuntime preserves the explicit-migration boundary even when a
// modern runtime with the same identifier also exists. It never modifies or
// adopts historical runtime directories.
func rejectLegacyRuntime(root string, runtime microruntime.Identifier) error {
	for _, directory := range []string{".micro-runtime", ".micro-runtimes"} {
		relative := filepath.Join(directory, runtime.String())
		_, err := os.Lstat(filepath.Join(root, relative))
		if err == nil {
			return fmt.Errorf("legacy runtime %s exists; migration is explicit, and this command does not modify it", relative)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}
