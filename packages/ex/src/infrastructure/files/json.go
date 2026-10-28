package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// maxDocumentBytes bounds every JSON document EX reads.
const maxDocumentBytes = 4 << 20

// Permissions for EX's private files and the directories that hold them.
const (
	privateFileMode      = 0o600
	privateDirectoryMode = 0o700
)

// ReadBounded reads a regular file reached without symlinks and fails rather
// than truncating when it exceeds limit bytes.
func ReadBounded(path string, limit int64) ([]byte, error) {
	if err := RequireNoSymlinks(filepath.Dir(path)); err != nil {
		return nil, err
	}

	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("expected a regular file: %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%s exceeds %d bytes", path, limit)
	}

	return data, nil
}

// ReadJSON reads one bounded JSON document. Unknown fields are allowed, so
// private records can evolve additively.
func ReadJSON(path string, value any) error {
	data, err := ReadBounded(path, maxDocumentBytes)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	return nil
}

// WriteJSON replaces a private document only after its complete content is
// synced, then syncs the directory, so a completed state transition never
// depends on a buffered directory entry. Existing links are refused.
func WriteJSON(path string, value any) error {
	directory := filepath.Dir(path)

	if err := RequireNoSymlinks(directory); err != nil {
		return err
	}

	if err := os.MkdirAll(directory, privateDirectoryMode); err != nil {
		return err
	}

	if info, err := os.Lstat(path); err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refusing to replace non-regular file %s", path)
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	temporary, err := writeTemporary(directory, append(data, '\n'))
	if err != nil {
		return err
	}
	defer os.Remove(temporary)

	if err := os.Rename(temporary, path); err != nil {
		return err
	}

	return syncDirectory(directory)
}

func writeTemporary(directory string, data []byte) (string, error) {
	file, err := os.CreateTemp(directory, ".state-*")
	if err != nil {
		return "", err
	}

	writeErr := file.Chmod(privateFileMode)
	if writeErr == nil {
		_, writeErr = file.Write(data)
	}

	if writeErr == nil {
		writeErr = file.Sync()
	}

	if err := errors.Join(writeErr, file.Close()); err != nil {
		_ = os.Remove(file.Name())

		return "", err
	}

	return file.Name(), nil
}

func syncDirectory(directory string) error {
	handle, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer handle.Close()

	return handle.Sync()
}
