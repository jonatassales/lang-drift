package configuration

import (
	"errors"
	"strings"
	"testing"
)

func TestParseEnvFile_Refuses(t *testing.T) {
	for _, test := range []struct {
		name    string
		content string
		reason  string
	}{
		{"lowercase name", "log_url=x\n", "uppercase variable name"},
		{"leading digit", "1PORT=x\n", "uppercase variable name"},
		{"export prefix", "export PORT=1\n", "uppercase variable name"},
		{"no assignment", "PORT\n", "NAME=value"},
		{"repeated name", "PORT=1\nPORT=2\n", "named twice"},
		{"control character", "PORT=1\x00\n", "plain text"},
		{"long value", "PORT=" + strings.Repeat("a", 4097) + "\n", "4 KiB"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseEnvFile(test.content)

			refusedForReason := err != nil &&
				strings.Contains(err.Error(), test.reason) &&
				errors.Is(err, ErrInvalidVariable)

			if !refusedForReason {
				t.Fatalf("ParseEnvFile err = %v, want %q", err, test.reason)
			}
		})
	}
}
