// Package files performs the local file operations EX trusts: path resolution
// that never follows symlinks, bounded reads, atomic private writes, and
// cross-process locks. It is not a sandbox against a hostile workspace; the
// consumer must not change the workspace while an operation runs.
package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CanonicalDirectory resolves path to an absolute directory with every
// symlink resolved, so two spellings of one directory compare equal.
func CanonicalDirectory(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(canonical)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", canonical)
	}

	return canonical, nil
}

// Within resolves relative below base, which may not exist yet, refusing
// symlinks in base's ancestors and in every existing component.
func Within(base, relative string) (string, error) {
	return resolve(base, relative, false)
}

// Existing resolves relative below boundary like Within, and also requires
// every component to exist.
func Existing(boundary, relative string) (string, error) {
	return resolve(boundary, relative, true)
}

func resolve(base, relative string, mustExist bool) (string, error) {
	if !filepath.IsLocal(relative) {
		return "", fmt.Errorf("path %q escapes its boundary", relative)
	}

	root, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	if err := RequireNoSymlinks(root); err != nil {
		return "", err
	}

	current := root

	for _, component := range strings.Split(filepath.Clean(relative), string(filepath.Separator)) {
		if component == "." {
			continue
		}

		current = filepath.Join(current, component)

		info, err := os.Lstat(current)

		switch {
		case err == nil && info.Mode()&os.ModeSymlink != 0:
			return "", fmt.Errorf("symlink path is not allowed: %s", current)

		case errors.Is(err, fs.ErrNotExist) && !mustExist:
			continue

		case err != nil:
			return "", fmt.Errorf("reference %s: %w", relative, err)
		}
	}

	return current, nil
}

// RequireNoSymlinks rejects path when it or any existing ancestor is a symlink.
func RequireNoSymlinks(path string) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	for current := absolute; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)

		switch {
		case err == nil && info.Mode()&os.ModeSymlink != 0:
			return fmt.Errorf("symlink path is not allowed: %s", current)

		case err != nil && !errors.Is(err, fs.ErrNotExist):
			return err
		}

		if current == filepath.Dir(current) {
			return nil
		}
	}
}
