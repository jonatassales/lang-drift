package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Directories that hold the documents a consumer reads: the folder is theirs.
const documentDirectoryMode = 0o755

// WriteDocument replaces one of the documents the LangDrift CLI writes for the consumer and a
// console to read, records and mirrors, with the given bytes and mode. The
// content is synced before the rename and the directory after it, as with
// private JSON, and missing directories are created readable. Existing links
// are refused.
func WriteDocument(path string, data []byte, mode fs.FileMode) error {
	directory := filepath.Dir(path)

	if err := RequireNoSymlinks(directory); err != nil {
		return err
	}

	if err := os.MkdirAll(directory, documentDirectoryMode); err != nil {
		return err
	}

	info, err := os.Lstat(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refusing to replace non-regular file %s", path)
	}

	temporary, err := writeTemporary(directory, data, mode)
	if err != nil {
		return err
	}
	defer os.Remove(temporary)

	if err := os.Rename(temporary, path); err != nil {
		return err
	}

	return syncDirectory(directory)
}
