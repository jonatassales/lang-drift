// Package session renders the events of a plan or a run as the website's
// terminal prints them: the same words, with the actor named agent and true
// figures, for people, or one JSON event per line for a console.
package session

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jonatassales/tao/packages/ex/src/application/session"
)

// Actor is who the lines say acted.
const Actor = "agent › "

// Dollars prints cents as dollars with cents, grouped by thousands.
func Dollars(cents int) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}

	return fmt.Sprintf("%s$%s.%02d", sign, group(cents/100), cents%100)
}

// Tokens prints thousands of tokens compactly: 3.5M, 250K, 6M.
func Tokens(thousands int) string {
	if thousands < 1000 {
		return strconv.Itoa(thousands) + "K"
	}

	millions := float64(thousands) / 1000

	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.1f", millions), "0"), ".") + "M"
}

// Number prints an integer grouped by thousands.
func Number(value int) string {
	if value < 0 {
		return "-" + group(-value)
	}

	return group(value)
}

func group(value int) string {
	digits := strconv.Itoa(value)
	if len(digits) <= 3 {
		return digits
	}

	var out strings.Builder
	head := len(digits) % 3
	if head > 0 {
		out.WriteString(digits[:head])
	}
	for index := head; index < len(digits); index += 3 {
		if out.Len() > 0 {
			out.WriteString(",")
		}
		out.WriteString(digits[index : index+3])
	}

	return out.String()
}

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
