package support_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"
)

func TestGlobals_Selection(t *testing.T) {
	t.Parallel()

	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	relative, err := filepath.Abs("nested/project")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		globals support.Globals
		want    application.ProjectSelection
	}{
		{name: "explicit project", globals: support.Globals{Project: "/workspace", Fallback: "/elsewhere"}, want: application.ProjectSelection{Directory: "/workspace", Explicit: true}},
		{name: "relative explicit project", globals: support.Globals{Project: "nested/project"}, want: application.ProjectSelection{Directory: relative, Explicit: true}},
		{name: "fallback directory", globals: support.Globals{Fallback: "/workspace/services"}, want: application.ProjectSelection{Directory: "/workspace/services"}},
		{name: "working directory", want: application.ProjectSelection{Directory: workingDirectory}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, err := tt.globals.Selection(); err != nil || got != tt.want {
				t.Fatalf("Selection() = %+v, %v; want %+v", got, err, tt.want)
			}
		})
	}
}
