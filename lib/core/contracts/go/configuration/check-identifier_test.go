package configuration

import (
	"errors"
	"strings"
	"testing"
)

func TestCheckIdentifier(t *testing.T) {
	for _, test := range []struct {
		text string
		ok   bool
	}{
		{"console-local", true},
		{"a", true},
		{strings.Repeat("a", 63), true},
		{strings.Repeat("a", 64), false},
		{"Console", false},
		{"-console", false},
		{"console-", false},
		{"console_local", false},
		{"", false},
	} {
		err := CheckIdentifier(test.text)
		wrongError := err != nil && !errors.Is(err, ErrInvalidIdentifier)
		unexpected := (err == nil) != test.ok || wrongError

		if unexpected {
			t.Errorf("CheckIdentifier(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}
