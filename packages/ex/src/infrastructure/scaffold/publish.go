// Package scaffold publishes new consumer-owned directories: generated
// runtime configuration and agent projects. It never merges into, or
// overwrites, an existing directory.
package scaffold

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/exphere/ex/src/infrastructure/files"
)

// Permissions of published, consumer-owned content.
const (
	directoryMode = 0o755
	fileMode      = 0o644
)

// Publish builds a scaffold in a sibling staging directory, then reserves the
// destination exclusively and publishes the staged entries without
// overwriting anything. It returns the published directory.
//
// Publication is not atomic. A handled failure removes only entries that
// still belong to this invocation and reports output it could not remove. An
// abrupt termination may leave staging or a partial directory behind; later
// calls refuse that destination rather than overwrite it.
func Publish(root, relative string, build func(stage string) error) (location string, resultErr error) {
	canonicalRoot, err := files.CanonicalDirectory(root)
	if err != nil {
		return "", err
	}

	parent, target, err := reserveParents(canonicalRoot, relative)
	if err != nil {
		return "", err
	}

	if _, err := os.Lstat(target); err == nil {
		return "", fmt.Errorf("destination already exists: %s", target)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	staging, err := os.MkdirTemp(parent, ".exphere-stage-")
	if err != nil {
		return "", err
	}

	defer func() {
		if cleanupErr := os.RemoveAll(staging); cleanupErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("retained scaffold staging directory %s: %w", staging, cleanupErr))
		}
	}()

	if err := build(staging); err != nil {
		return "", err
	}

	if err := os.Chmod(staging, directoryMode); err != nil {
		return "", err
	}

	if err := publish(staging, target); err != nil {
		return "", err
	}

	return target, nil
}

// reserveParents creates the destination's parent directories, refusing any
// that is not a real directory, and returns the parent and the destination.
func reserveParents(root, relative string) (parent, target string, err error) {
	if !filepath.IsLocal(relative) || filepath.Clean(relative) != relative || strings.Contains(relative, `\`) {
		return "", "", fmt.Errorf("scaffold destination %q must be a clean relative path", relative)
	}

	segments := strings.Split(filepath.ToSlash(relative), "/")
	parent = root

	for _, segment := range segments[:len(segments)-1] {
		parent = filepath.Join(parent, segment)

		if err := os.Mkdir(parent, directoryMode); err != nil && !errors.Is(err, fs.ErrExist) {
			return "", "", err
		}

		info, err := os.Lstat(parent)
		if err != nil {
			return "", "", err
		}

		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return "", "", fmt.Errorf("scaffold parent must be a real directory: %s", parent)
		}
	}

	return parent, filepath.Join(parent, segments[len(segments)-1]), nil
}

// publish reserves target and links every staged file into it. Staging
// shares target's filesystem, so a link refuses an existing destination and
// keeps the file's permissions; removing staging later removes only its link.
func publish(staging, target string) (resultErr error) {
	if err := os.Mkdir(target, directoryMode); err != nil {
		return fmt.Errorf("reserve scaffold destination: %w", err)
	}

	targetInfo, err := os.Lstat(target)
	if err != nil {
		return fmt.Errorf("inspect reserved scaffold destination %s: %w", target, err)
	}

	owned := ownedEntries{{path: target, info: targetInfo, mode: directoryMode}}

	defer func() {
		if resultErr != nil {
			resultErr = errors.Join(resultErr, owned.remove())
		}
	}()

	walkErr := filepath.WalkDir(staging, func(source string, _ fs.DirEntry, walkErr error) error {
		if walkErr != nil || source == staging {
			return walkErr
		}

		relative, err := filepath.Rel(staging, source)
		if err != nil {
			return err
		}

		entry, err := owned.publishEntry(source, filepath.Join(target, relative))
		if err != nil {
			return err
		}

		owned = append(owned, entry)

		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("publish scaffold from %s: %w", staging, walkErr)
	}

	return owned.finalizeDirectories()
}

// ownedEntry is one published path and the identity it had when published.
type ownedEntry struct {
	path string
	info fs.FileInfo
	mode fs.FileMode
}

// ownedEntries are published paths in creation order: parents before children.
type ownedEntries []ownedEntry

func (owned ownedEntries) publishEntry(source, destination string) (ownedEntry, error) {
	info, err := os.Lstat(source)
	if err != nil {
		return ownedEntry{}, err
	}

	if !info.IsDir() && !info.Mode().IsRegular() {
		return ownedEntry{}, fmt.Errorf("unsupported staged scaffold entry: %s", source)
	}

	if err := owned.checkAncestors(destination); err != nil {
		return ownedEntry{}, err
	}

	if !info.IsDir() {
		if err := os.Link(source, destination); err != nil {
			return ownedEntry{}, fmt.Errorf("publish scaffold file %s: %w", destination, err)
		}

		return ownedEntry{path: destination, info: info}, nil
	}

	if err := os.Mkdir(destination, directoryMode); err != nil {
		return ownedEntry{}, fmt.Errorf("reserve scaffold directory %s: %w", destination, err)
	}

	published, err := os.Lstat(destination)
	if err != nil {
		return ownedEntry{}, err
	}

	return ownedEntry{path: destination, info: published, mode: info.Mode().Perm()}, nil
}

// finalizeDirectories sets directory permissions after their children are
// published, so a read-only staged directory can still be populated.
func (owned ownedEntries) finalizeDirectories() error {
	for index := len(owned) - 1; index >= 0; index-- {
		entry := owned[index]
		if !entry.info.IsDir() {
			continue
		}

		if err := owned[:index].checkAncestors(entry.path); err != nil {
			return err
		}

		if err := entry.checkIdentity(); err != nil {
			return err
		}

		if err := os.Chmod(entry.path, entry.mode); err != nil {
			return fmt.Errorf("set scaffold directory permissions %s: %w", entry.path, err)
		}
	}

	return nil
}

// remove deletes published entries, children before parents. Remove rather
// than RemoveAll keeps files someone else added to an owned directory.
func (owned ownedEntries) remove() error {
	var result error

	for index := len(owned) - 1; index >= 0; index-- {
		entry := owned[index]

		if err := owned[:index].checkAncestors(entry.path); err != nil {
			return errors.Join(result, fmt.Errorf("retained scaffold output: %w", err))
		}

		if err := entry.checkIdentity(); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				result = errors.Join(result, fmt.Errorf("retained scaffold output: %w", err))
			}

			continue
		}

		if err := os.Remove(entry.path); err != nil {
			result = errors.Join(result, fmt.Errorf("retained scaffold output %s: %w", entry.path, err))
		}
	}

	return result
}

// checkAncestors refuses a replaced parent, including a symlink swapped in,
// even when a descendant keeps its identity through another path.
func (owned ownedEntries) checkAncestors(path string) error {
	for _, parent := range owned {
		if parent.info.IsDir() && strings.HasPrefix(path, parent.path+string(filepath.Separator)) {
			if err := parent.checkIdentity(); err != nil {
				return fmt.Errorf("check scaffold parent: %w", err)
			}
		}
	}

	return nil
}

func (entry ownedEntry) checkIdentity() error {
	info, err := os.Lstat(entry.path)
	if err != nil {
		return fmt.Errorf("inspect scaffold ownership %s: %w", entry.path, err)
	}

	if !os.SameFile(entry.info, info) {
		return fmt.Errorf("scaffold ownership changed: %s", entry.path)
	}

	return nil
}

// WriteFile creates a new file beneath a staging directory.
func WriteFile(stage, relative string, data []byte) error {
	if !fs.ValidPath(relative) || relative == "." || strings.Contains(relative, `\`) {
		return fmt.Errorf("invalid scaffold file path %q", relative)
	}

	path := filepath.Join(stage, filepath.FromSlash(relative))

	if err := os.MkdirAll(filepath.Dir(path), directoryMode); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fileMode)
	if err != nil {
		return err
	}

	_, writeErr := file.Write(data)

	return errors.Join(writeErr, file.Close())
}

// WriteJSONFile creates a new, indented JSON file beneath a staging directory.
func WriteJSONFile(stage, relative string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return WriteFile(stage, relative, append(data, '\n'))
}
