package main

import (
	"io/fs"
	"path/filepath"
	"testing"
)

// filesUnder lists every file below a folder.
func filesUnder(t *testing.T, root string) []string {
	t.Helper()

	var found []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if !entry.IsDir() {
			found = append(found, path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return found
}
