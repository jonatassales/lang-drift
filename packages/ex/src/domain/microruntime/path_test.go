package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

// relativePath parses text that the test knows to be valid.
func relativePath(t testing.TB, text string) microruntime.RelativePath {
	t.Helper()

	parsed, err := microruntime.ParseRelativePath(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

func TestParseRelativePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{name: "single segment", text: "terraform"},
		{name: "nested", text: "services/web"},
		{name: "hidden segment", text: "units/.generated"},
		{name: "empty", text: "", wantErr: true},
		{name: "blank", text: "   ", wantErr: true},
		{name: "current directory", text: ".", wantErr: true},
		{name: "parent directory", text: "..", wantErr: true},
		{name: "escapes upward", text: "../web", wantErr: true},
		{name: "unclean traversal", text: "services/../web", wantErr: true},
		{name: "unclean current segment", text: "services/./web", wantErr: true},
		{name: "trailing slash", text: "services/", wantErr: true},
		{name: "doubled slash", text: "services//web", wantErr: true},
		{name: "absolute", text: "/services/web", wantErr: true},
		{name: "backslash", text: `services\web`, wantErr: true},
		{name: "control character", text: "services/\tweb", wantErr: true},
		{name: "NUL byte", text: "services/web\x00", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := microruntime.ParseRelativePath(tt.text)

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrInvalidPath) || !got.IsZero() {
					t.Fatalf("ParseRelativePath(%q) = %q, %v; want the zero value and ErrInvalidPath", tt.text, got, err)
				}

				return
			}

			if err != nil || got.String() != tt.text {
				t.Fatalf("ParseRelativePath(%q) = %q, %v", tt.text, got, err)
			}
		})
	}
}

func TestRelativePath_Accessors(t *testing.T) {
	t.Parallel()

	path := relativePath(t, "hooks/repair.sh")

	if path.Base() != "repair.sh" {
		t.Errorf("Base() = %q, want repair.sh", path.Base())
	}

	if !path.HasSuffix(".sh") || path.HasSuffix(".json") {
		t.Errorf("HasSuffix() does not report the .sh suffix correctly for %q", path)
	}
}

func TestRelativePath_IsWithin(t *testing.T) {
	t.Parallel()

	runtimes := relativePath(t, ".exphere/runtimes")

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "the directory itself", path: ".exphere/runtimes", want: true},
		{name: "a descendant", path: ".exphere/runtimes/demo/units", want: true},
		{name: "a sibling sharing a prefix", path: ".exphere/runtimes-old", want: false},
		{name: "the parent", path: ".exphere", want: false},
		{name: "an unrelated path", path: "services/web", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := relativePath(t, tt.path).IsWithin(runtimes); got != tt.want {
				t.Fatalf("%q.IsWithin(%q) = %t, want %t", tt.path, runtimes, got, tt.want)
			}
		})
	}
}
