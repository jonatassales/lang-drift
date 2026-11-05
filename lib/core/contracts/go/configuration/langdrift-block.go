package configuration

import "strings"

// maximumMarkdownBytes bounds the document a fence is read from.
const maximumMarkdownBytes = 1 << 20

// langdriftBlock is the body of the one langdrift JSON fence in a bounded
// Markdown document.
func langdriftBlock(markdown string) (string, error) {
	if len(markdown) > maximumMarkdownBytes {
		return "", invalid("Markdown exceeds 1 MiB")
	}

	const opening = "```langdrift\n"
	if strings.Count(markdown, opening) != 1 {
		return "", invalid("Markdown requires exactly one langdrift JSON fence")
	}

	after := strings.SplitN(markdown, opening, 2)[1]

	end := strings.Index(after, "\n```")
	if end < 0 {
		return "", invalid("configuration fence is not closed")
	}

	return after[:end], nil
}
