package project_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/exphere/ex/src/domain/project"
)

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{name: "lowercase hexadecimal", text: "0123456789abcdef0123456789abcdef"},
		{name: "empty", text: "", wantErr: true},
		{name: "too short", text: strings.Repeat("a", 31), wantErr: true},
		{name: "too long", text: strings.Repeat("a", 33), wantErr: true},
		{name: "uppercase", text: strings.Repeat("A", 32), wantErr: true},
		{name: "not hexadecimal", text: strings.Repeat("g", 32), wantErr: true},
		{name: "leading space", text: " " + strings.Repeat("a", 31), wantErr: true},
		{name: "trailing newline", text: strings.Repeat("a", 32) + "\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := project.ParseID(tt.text)

			if tt.wantErr {
				if !errors.Is(err, project.ErrInvalidID) {
					t.Fatalf("ParseID(%q) error = %v, want ErrInvalidID", tt.text, err)
				}

				if id != (project.ID{}) {
					t.Fatalf("ParseID(%q) = %q on error, want the zero ID", tt.text, id)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseID(%q) error = %v", tt.text, err)
			}

			if id.String() != tt.text {
				t.Fatalf("ParseID(%q).String() = %q", tt.text, id)
			}
		})
	}
}

func TestProject_ClusterName(t *testing.T) {
	t.Parallel()

	id, err := project.ParseID("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}

	got := project.Project{ID: id, Root: "/workspace"}.ClusterName()
	if got != "ex-0123456789ab" {
		t.Fatalf("ClusterName() = %q, want ex-0123456789ab", got)
	}
}
