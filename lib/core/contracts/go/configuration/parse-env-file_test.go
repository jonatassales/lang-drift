package configuration

import (
	"errors"
	"strings"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	variables, err := ParseEnvFile("# the log\nLOG_URL=http://log:3500\n\nPORT=3402\nEMPTY=\nQUERY=a=b\n")
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"LOG_URL": "http://log:3500",
		"PORT":    "3402",
		"EMPTY":   "",
		"QUERY":   "a=b",
	}
	if len(variables) != len(want) {
		t.Fatalf("variables = %v, want %v", variables, want)
	}

	for name, value := range want {
		if variables[name] != value {
			t.Fatalf("%s = %q, want %q", name, variables[name], value)
		}
	}
}

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

			if err == nil || !strings.Contains(err.Error(), test.reason) || !errors.Is(err, ErrInvalidVariable) {
				t.Fatalf("ParseEnvFile err = %v, want %q", err, test.reason)
			}
		})
	}
}
