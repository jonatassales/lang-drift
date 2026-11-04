package collection

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// Syntax is how a link was written.
type Syntax string

// The link syntaxes, as the records name them.
const (
	NamedInCLI   Syntax = "named in the CLI"
	MarkdownLink Syntax = "markdown link"
	WikiLink     Syntax = "wiki link"
	FieldLine    Syntax = "field line"
)

// Link is one reference a document makes to another file.
type Link struct {
	Target string
	Syntax Syntax
}

var (
	markdownLink = regexp.MustCompile(`\[[^\]]*\]\(([^)\s]+)\)`)
	wikiLink     = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]*)?\]\]`)
	fieldLine    = regexp.MustCompile(`(?m)^\*\*[^*\n]+:\*\* +(\S+)\s*$`)
	scheme       = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:`)
)

// Links finds the links of a Markdown document, in document order: markdown
// links and wiki links to files, and field lines whose value is a path. A
// link with a scheme, a fragment-only link and a field value without a slash
// are not links to files.
func Links(markdown string) []Link {
	type found struct {
		offset int
		link   Link
	}

	var links []found

	for _, match := range markdownLink.FindAllStringSubmatchIndex(markdown, -1) {
		target := markdown[match[2]:match[3]]
		if isFileTarget(target) {
			links = append(links, found{match[0], Link{Target: target, Syntax: MarkdownLink}})
		}
	}

	for _, match := range wikiLink.FindAllStringSubmatchIndex(markdown, -1) {
		target := markdown[match[2]:match[3]]
		if isFileTarget(target) {
			links = append(links, found{match[0], Link{Target: target, Syntax: WikiLink}})
		}
	}

	for _, match := range fieldLine.FindAllStringSubmatchIndex(markdown, -1) {
		target := markdown[match[2]:match[3]]
		if isFileTarget(target) && strings.Contains(target, "/") {
			links = append(links, found{match[0], Link{Target: target, Syntax: FieldLine}})
		}
	}

	sort.SliceStable(links, func(left, right int) bool { return links[left].offset < links[right].offset })

	ordered := make([]Link, 0, len(links))
	for _, item := range links {
		ordered = append(ordered, item.link)
	}

	return ordered
}

func isFileTarget(target string) bool {
	return target != "" && !strings.HasPrefix(target, "#") && !scheme.MatchString(target)
}

// Candidates are the project-relative paths a link may name, in the order
// to try them: a target that starts with a slash names a path from the
// project root; a field line's value is project-relative by convention;
// any other target is relative to the linking file, with the project root
// as a second reading. Fragments are stripped.
func Candidates(link Link, linkingFile string) []string {
	target, _, _ := strings.Cut(link.Target, "#")

	switch {
	case strings.HasPrefix(target, "/"):
		return []string{path.Clean(strings.TrimPrefix(target, "/"))}

	case link.Syntax == FieldLine:
		return []string{path.Clean(target)}
	}

	relative := path.Clean(path.Join(path.Dir(linkingFile), target))
	fromRoot := path.Clean(target)
	if relative == fromRoot {
		return []string{relative}
	}

	return []string{relative, fromRoot}
}

// Escapes reports whether a candidate leaves the project.
func Escapes(candidate string) bool {
	return candidate == ".." || strings.HasPrefix(candidate, "../") || path.IsAbs(candidate)
}
