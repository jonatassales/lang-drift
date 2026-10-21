package connected

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
)

func hasGitRoot(root string) (bool, error) {
	for directory := root; ; directory = filepath.Dir(directory) {
		_, err := os.Lstat(filepath.Join(directory, ".git"))
		if err == nil {
			return true, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
		if directory == filepath.Dir(directory) {
			return false, nil
		}
	}
}

func excludedSource(name string) bool {
	for _, part := range strings.Split(filepath.ToSlash(name), "/") {
		if part == ".env" || strings.HasPrefix(part, ".env.") {
			return true
		}
		switch part {
		case ".git", ".langdrift", "node_modules", "vendor", ".next", "dist", "build", "target", ".cache", ".venv", "venv", "__pycache__":
			return true
		}
	}
	return false
}

type ignoreRule struct {
	pattern *regexp.Regexp
	include bool
}

// readIgnore supports directory-scoped Git ignore globs without requiring a
// Git installation for plain folders. Invalid patterns fail discovery closed.
func readIgnore(root, relative string, rules []ignoreRule) ([]ignoreRule, error) {
	path, err := files.Within(root, filepath.Join(relative, ".gitignore"))
	if err != nil {
		return nil, err
	}
	data, err := files.ReadBounded(path, 1<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return rules, nil
	}
	if err != nil {
		return nil, err
	}
	result := append([]ignoreRule{}, rules...)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, " \r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Match the console discovery subset; never literalize unsupported
		// classes or escapes and accidentally include ignored private files.
		if strings.ContainsAny(line, `[]\`) {
			return nil, errors.New("unsupported ignore pattern: character classes and escapes")
		}
		include := strings.HasPrefix(line, "!")
		if include {
			line = strings.TrimPrefix(line, "!")
		}
		anchored := strings.HasPrefix(line, "/") || strings.Contains(strings.TrimSuffix(line, "/"), "/")
		line = strings.Trim(line, "/")
		pattern := regexp.QuoteMeta(line)
		pattern = strings.ReplaceAll(pattern, `\*\*/`, `(?:.*/)?`)
		pattern = strings.ReplaceAll(pattern, `\*\*`, `.*`)
		pattern = strings.ReplaceAll(pattern, `\*`, `[^/]*`)
		pattern = strings.ReplaceAll(pattern, `\?`, `[^/]`)
		prefix := ""
		if relative != "." {
			prefix = regexp.QuoteMeta(filepath.ToSlash(relative)) + "/"
		}
		if !anchored {
			prefix += "(?:.*/)?"
		}
		compiled, err := regexp.Compile("^" + prefix + pattern + "(?:/.*)?$")
		if err != nil {
			return nil, err
		}
		result = append(result, ignoreRule{compiled, include})
	}
	return result, nil
}

func walkLocal(ctx context.Context, root string) ([]string, error) {
	var names []string
	count := 0
	var visit func(string, []ignoreRule) error
	visit = func(relative string, inherited []ignoreRule) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		rules, err := readIgnore(root, relative, inherited)
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(filepath.Join(root, relative))
		if err != nil {
			return err
		}
		for _, entry := range entries {
			count++
			if count > 20000 {
				return errors.New("local source exceeds 20,000 entries")
			}
			name := filepath.Join(relative, entry.Name())
			if excludedSource(name) || entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			ignored := false
			for _, rule := range rules {
				if rule.pattern.MatchString(filepath.ToSlash(name)) {
					ignored = !rule.include
				}
			}
			if ignored {
				continue
			}
			if entry.IsDir() {
				if err := visit(name, rules); err != nil {
					return err
				}
				continue
			}
			if entry.Type().IsRegular() {
				names = append(names, name)
			}
		}
		return nil
	}
	err := visit(".", nil)
	return names, err
}

func (native Native) localBuildTree(ctx context.Context, owner project.Project) (string, func(), string, error) {
	git, err := hasGitRoot(owner.Root)
	if err != nil {
		return "", nil, "", err
	}
	var names []string
	if git {
		result, readErr := native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", "."}, Dir: owner.Root})
		if readErr != nil {
			return "", nil, "", readErr
		}
		names = strings.Split(result.Stdout, "\x00")
	} else {
		names, err = walkLocal(ctx, owner.Root)
		if err != nil {
			return "", nil, "", err
		}
	}
	sort.Strings(names)
	filtered := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" || seen[name] || excludedSource(name) {
			continue
		}
		seen[name] = true
		path, pathErr := files.Within(owner.Root, name)
		if pathErr != nil {
			continue
		} // Symlinks never enter a local build context.
		info, statErr := os.Lstat(path)
		if errors.Is(statErr, fs.ErrNotExist) {
			continue
		} // Include local deletions.
		if statErr != nil {
			return "", nil, "", statErr
		}
		if info.Mode().IsRegular() {
			filtered = append(filtered, name)
		}
	}
	root, cleanup, err := copyBuildTree(ctx, owner, filtered)
	if err != nil {
		return "", nil, "", err
	}
	hash := sha256.New()
	for _, name := range filtered {
		data, err := files.ReadBounded(filepath.Join(root, name), 8<<20)
		if err != nil {
			cleanup()
			return "", nil, "", err
		}
		_, _ = hash.Write([]byte(name + "\x00" + digest(data) + "\x00"))
	}
	return root, cleanup, hex.EncodeToString(hash.Sum(nil)), nil
}
