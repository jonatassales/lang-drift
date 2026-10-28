// Package collection is the filesystem side of the deterministic shell: it
// follows a ticket's links breadth-first under the limits, refuses secrets,
// symlinks and paths outside the project, and records every file it reaches
// by type, size and hash, never by content.
package collection

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// Collector follows links on disk.
type Collector struct{}

type pending struct {
	candidates []string
	depth      int
	linkedFrom string
	syntax     collection.Syntax
}

// Collect follows the ticket's links under the limits. The ticket path is
// project-relative.
func (Collector) Collect(owner project.Project, ticketPath string, limits collection.Limits) (collection.Result, error) {
	ticketPath = filepath.ToSlash(filepath.Clean(ticketPath))

	queue := []pending{{candidates: []string{ticketPath}, depth: 0, syntax: collection.NamedInCLI}}
	visited := map[string]bool{}
	hashes := map[string]string{}
	result := collection.Result{}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		source, content, found := resolve(owner, next)
		if visited[source.Path] {
			continue
		}
		visited[source.Path] = true

		if found && source.State == "" {
			if result.Files >= limits.MaxFiles || result.Bytes+int64(len(content)) > limits.MaxBytes {
				// Beyond the limits: the file is known, not read.
				source.State = collection.Unresolved
				source.MediaType, source.Bytes, source.SHA256 = "", 0, ""
				result.Sources = append(result.Sources, source)

				continue
			}

			describe(&source, content)
			result.Files++
			result.Bytes += source.Bytes
			result.Depth = max(result.Depth, source.Depth)
			hashes[source.Path] = source.SHA256
		}

		result.Sources = append(result.Sources, source)

		if source.State != collection.Read || source.MediaType != collection.MarkdownType || source.Depth >= limits.MaxDepth {
			continue
		}

		for _, link := range collection.Links(string(content)) {
			queue = append(queue, pending{
				candidates: collection.Candidates(link, source.Path),
				depth:      source.Depth + 1,
				linkedFrom: source.Path,
				syntax:     link.Syntax,
			})
		}
	}

	result.ContextSHA256 = collection.ContextHash(hashes)

	return result, nil
}

// resolve finds the first candidate that exists, or records why none does.
// The returned state is empty for a file that was read and still has to be
// described, so the caller can apply the limits first.
func resolve(owner project.Project, next pending) (collection.Source, []byte, bool) {
	source := collection.Source{Depth: next.depth, LinkedFrom: next.linkedFrom, LinkSyntax: next.syntax, Path: next.candidates[0]}

	for _, candidate := range next.candidates {
		source.Path = candidate

		if collection.Escapes(candidate) {
			source.State = collection.Refused

			return source, nil, false
		}

		if collection.IsSecret(candidate) {
			source.State = collection.Refused

			return source, nil, false
		}

		absolute, err := files.Existing(owner.Root, filepath.FromSlash(candidate))
		if err != nil && strings.Contains(err.Error(), "symlink") {
			source.State = collection.Refused

			return source, nil, false
		}
		if errors.Is(err, fs.ErrNotExist) || err != nil {
			continue
		}

		info, err := os.Lstat(absolute)
		if err != nil || !info.Mode().IsRegular() {
			continue
		}

		content, err := files.ReadBounded(absolute, 64<<20)
		if err != nil {
			continue
		}

		return source, content, true
	}

	source.Path = next.candidates[0]
	source.State = collection.Unresolved

	return source, nil, false
}

// describe records a read file by type, size and hash.
func describe(source *collection.Source, content []byte) {
	source.Bytes = int64(len(content))
	source.SHA256 = collection.SHA256(content)
	source.Content = content

	if collection.IsBinary(content) {
		source.State = collection.Binary
		source.MediaType = collection.BinaryType
		source.Content = nil

		return
	}

	source.State = collection.Read
	source.MediaType = collection.MediaType(source.Path)
}
