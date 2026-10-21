package microruntime_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

// identifier parses text that the test knows to be valid.
func identifier(t testing.TB, text string) microruntime.Identifier {
	t.Helper()

	parsed, err := microruntime.ParseIdentifier(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

func TestParseIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{name: "single character", text: "a"},
		{name: "digits and hyphens", text: "web-2"},
		{name: "maximum length", text: strings.Repeat("a", 63)},
		{name: "empty", text: "", wantErr: true},
		{name: "too long", text: strings.Repeat("a", 64), wantErr: true},
		{name: "uppercase", text: "Web", wantErr: true},
		{name: "leading hyphen", text: "-web", wantErr: true},
		{name: "trailing hyphen", text: "web-", wantErr: true},
		{name: "underscore", text: "web_api", wantErr: true},
		{name: "dot", text: "web.api", wantErr: true},
		{name: "path separator", text: "web/api", wantErr: true},
		{name: "trailing newline", text: "web\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := microruntime.ParseIdentifier(tt.text)

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrInvalidIdentifier) || !got.IsZero() {
					t.Fatalf("ParseIdentifier(%q) = %q, %v; want the zero value and ErrInvalidIdentifier", tt.text, got, err)
				}

				return
			}

			if err != nil || got.String() != tt.text || got.IsZero() {
				t.Fatalf("ParseIdentifier(%q) = %q, %v", tt.text, got, err)
			}
		})
	}
}
