package session

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// The terminal's tones: a comment dimmed, a success line green, a warning
// yellow, and the reset that ends each.
const (
	dimmed = "\x1b[2m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	reset  = "\x1b[0m"
)

// tones are the kinds a terminal colors; every other kind keeps the
// terminal's own color.
var tones = map[session.Kind]string{
	session.Comment: dimmed,
	session.Success: green,
	session.Warning: yellow,
}

// toned wraps a line in its kind's tone, when the kind has one.
func toned(kind session.Kind, line string) string {
	tone, colored := tones[kind]
	if !colored {
		return line
	}

	return tone + line + reset
}
