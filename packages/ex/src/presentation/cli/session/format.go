// Package session renders the events of a plan or a run as the website's
// terminal prints them: the same words, with the actor named agent and true
// figures, for people, or one JSON event per line for a console.
package session

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jonatassales/tao/packages/ex/src/application/session"
)

// Actor is who the lines say acted.
const Actor = "agent › "

// pad right-pads text with spaces to width, as the website aligns its columns.
func pad(text string, width int) string {
	columns := utf8.RuneCountInString(text)
	if columns >= width {
		return text
	}

	return text + strings.Repeat(" ", width-columns)
}

func text(data session.Data, key string) string {
	switch value := data[key].(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	case nil:
		return ""
	}

	return fmt.Sprint(data[key])
}

func number(data session.Data, key string) int {
	switch value := data[key].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	}

	return 0
}

func list(data session.Data, key string) []string {
	switch value := data[key].(type) {
	case []string:
		return value
	case []any:
		items := make([]string, 0, len(value))
		for _, item := range value {
			items = append(items, fmt.Sprint(item))
		}

		return items
	}

	return nil
}

func joined(data session.Data, key string) string {
	return strings.Join(list(data, key), " · ")
}
