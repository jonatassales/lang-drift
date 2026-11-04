package manifest

import "strings"

// instruction is one Dockerfile instruction: its keyword in uppercase and
// its arguments, continuation lines joined.
type instruction struct {
	keyword   string
	arguments string
}

// instructions reads a Dockerfile's instructions in order. Comments and
// blank lines are skipped and a line ending in a backslash continues on
// the next, as Docker reads them.
func instructions(dockerfile []byte) []instruction {
	var found []instruction
	var pending strings.Builder

	for _, line := range strings.Split(string(dockerfile), "\n") {
		line = strings.TrimRight(line, " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if continued, continues := strings.CutSuffix(line, `\`); continues {
			pending.WriteString(continued + " ")

			continue
		}

		pending.WriteString(line)
		text := strings.TrimSpace(pending.String())
		pending.Reset()

		keyword, arguments := cutWord(text)
		found = append(found, instruction{
			keyword:   strings.ToUpper(keyword),
			arguments: arguments,
		})
	}

	return found
}
