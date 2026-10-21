package configuration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/domain/project"
)

func TestFiles_CheckReferences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, owner project.Project, location string)
		wantErr bool
	}{
		{name: "generated configuration", arrange: func(*testing.T, project.Project, string) {}},
		{
			name: "broken manifests that diagnostics must still tolerate",
			arrange: func(t *testing.T, _ project.Project, location string) {
				write(t, filepath.Join(location, "units", "web", "kubernetes", "service.json"), "{broken")
			},
		},
		{
			name: "missing component",
			arrange: func(t *testing.T, owner project.Project, _ string) {
				if err := os.RemoveAll(filepath.Join(owner.Root, "services")); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
		},
		{
			name: "Terraform directory without configuration",
			arrange: func(t *testing.T, _ project.Project, location string) {
				if err := os.Remove(filepath.Join(location, "terraform", "main.tf.json")); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
		},
		{
			name: "manifests reference that is not a manifest file",
			arrange: func(t *testing.T, _ project.Project, location string) {
				directory := filepath.Join(location, "units", "web", "kubernetes")
				if err := os.RemoveAll(directory); err != nil {
					t.Fatal(err)
				}
				write(t, directory, "not a directory or manifest")
			},
			wantErr: true,
		},
		{
			name: "missing manifests",
			arrange: func(t *testing.T, _ project.Project, location string) {
				if err := os.RemoveAll(filepath.Join(location, "units", "batch")); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner, location := generatedRuntime(t)
			tt.arrange(t, owner, location)

			descriptor := testBlueprint(t).Descriptor(identifier(t, "demo"))
			err := Files{}.CheckReferences(owner, descriptor)

			if tt.wantErr != (err != nil) {
				t.Fatalf("CheckReferences() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestFiles_ResolveHook(t *testing.T) {
	t.Parallel()

	owner, location := generatedRuntime(t)
	write(t, filepath.Join(location, "hooks", "seed.sh"), "#!/bin/sh\n")
	write(t, filepath.Join(location, "hooks", "notes.txt"), "")
	if err := os.MkdirAll(filepath.Join(location, "hooks", "directory.sh"), 0o755); err != nil {
		t.Fatal(err)
	}

	elsewhere := filepath.Join(tempDir(t), "linked.sh")
	write(t, elsewhere, "#!/bin/sh\n")
	if err := os.Symlink(elsewhere, filepath.Join(location, "hooks", "linked.sh")); err != nil {
		t.Fatal(err)
	}

	script, err := Files{}.ResolveHook(owner, identifier(t, "demo"), relativePath(t, "hooks/seed.sh"))
	if err != nil || script.Path != filepath.Join(location, "hooks", "seed.sh") || script.Directory != location {
		t.Fatalf("ResolveHook() = %+v, %v; want the script and the runtime directory", script, err)
	}

	for _, refused := range []string{"hooks/notes.txt", "hooks/directory.sh", "hooks/linked.sh", "hooks/missing.sh"} {
		if _, err := (Files{}).ResolveHook(owner, identifier(t, "demo"), relativePath(t, refused)); err == nil {
			t.Errorf("ResolveHook(%s) succeeded", refused)
		}
	}

	if _, err := (Files{}).ResolveHook(owner, identifier(t, "other"), relativePath(t, "hooks/seed.sh")); err == nil {
		t.Error("ResolveHook() resolved a hook of an uninitialized runtime")
	}
}
