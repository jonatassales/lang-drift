package configuration

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func loadAndSnapshot(t *testing.T, owner project.Project) (microruntime.Descriptor, application.Snapshot, error) {
	t.Helper()

	descriptor, err := Files{}.Load(owner, identifier(t, "demo"))
	if err != nil {
		return microruntime.Descriptor{}, application.Snapshot{}, err
	}

	snapshot, err := Files{}.Snapshot(owner, descriptor)

	return descriptor, snapshot, err
}

func TestFiles_LoadMapsTheDescriptor(t *testing.T) {
	t.Parallel()

	owner, _ := generatedRuntime(t)

	descriptor, err := Files{}.Load(owner, identifier(t, "demo"))
	if err != nil {
		t.Fatal(err)
	}

	want := testBlueprint(t).Descriptor(identifier(t, "demo"))
	if descriptor.ID != want.ID || descriptor.Owner != want.Owner || descriptor.Target != want.Target ||
		descriptor.TerraformDirectory != want.TerraformDirectory || len(descriptor.Units) != 2 ||
		descriptor.Units[identifier(t, "web")] != want.Units[identifier(t, "web")] {
		t.Fatalf("Load() = %+v, want %+v", descriptor, want)
	}
}

func TestFiles_LoadRefuses(t *testing.T) {
	t.Parallel()

	t.Run("uninitialized runtime", func(t *testing.T) {
		t.Parallel()

		if _, err := (Files{}).Load(testProject(t), identifier(t, "demo")); !errors.Is(err, application.ErrNotFound) {
			t.Fatalf("Load() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("legacy runtime with the same ID", func(t *testing.T) {
		t.Parallel()

		owner, _ := generatedRuntime(t)
		write(t, filepath.Join(owner.Root, ".micro-runtime", "demo"), "")

		if _, err := (Files{}).Load(owner, identifier(t, "demo")); err == nil || !strings.Contains(err.Error(), "legacy runtime") {
			t.Fatalf("Load() error = %v, want the legacy runtime refusal", err)
		}
	})

	t.Run("symlinked descriptor", func(t *testing.T) {
		t.Parallel()

		owner, location := generatedRuntime(t)
		descriptor := filepath.Join(location, "schema.json")
		moved := filepath.Join(tempDir(t), "schema.json")

		if err := os.Rename(descriptor, moved); err != nil {
			t.Fatal(err)
		}

		if err := os.Symlink(moved, descriptor); err != nil {
			t.Fatal(err)
		}

		if _, err := (Files{}).Load(owner, identifier(t, "demo")); err == nil {
			t.Fatal("Load() followed a symlinked descriptor")
		}
	})

	t.Run("descriptor that does not decode strictly", func(t *testing.T) {
		t.Parallel()

		owner, location := generatedRuntime(t)
		write(t, filepath.Join(location, "schema.json"), `{"version": 2, "version": 2}`)

		if _, err := (Files{}).Load(owner, identifier(t, "demo")); !errors.Is(err, microruntime.ErrInvalidDescriptor) {
			t.Fatalf("Load() error = %v, want ErrInvalidDescriptor", err)
		}
	})
}

func TestFiles_Snapshot(t *testing.T) {
	t.Parallel()

	owner, location := generatedRuntime(t)

	_, first, err := loadAndSnapshot(t, owner)
	if err != nil {
		t.Fatal(err)
	}

	var keys []string
	for _, manifest := range first.Manifests {
		keys = append(keys, manifest.Unit.String()+" "+manifest.Key())
	}

	if want := []string{"batch Job/batch", "web Service/web", "web Deployment/web"}; !slices.Equal(keys, want) {
		t.Errorf("manifests = %v, want units in order and files in lexical order: %v", keys, want)
	}

	if first.PreparationDirectory != filepath.Join(location, "terraform") {
		t.Errorf("PreparationDirectory = %q", first.PreparationDirectory)
	}

	if len(first.Hash) != 64 || len(first.CompatibleHashes) != 1 || first.CompatibleHashes[0] == first.Hash {
		t.Errorf("Hash = %q, CompatibleHashes = %v; want a digest and one distinct earlier encoding", first.Hash, first.CompatibleHashes)
	}

	_, second, err := loadAndSnapshot(t, owner)
	if err != nil || second.Hash != first.Hash || !slices.Equal(second.CompatibleHashes, first.CompatibleHashes) {
		t.Errorf("a second snapshot of unchanged files = %q, %v, %v; want the same fingerprints", second.Hash, second.CompatibleHashes, err)
	}
}

func TestFiles_SnapshotFingerprintTracksReferencedFilesOnly(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		edit        func(t *testing.T, location string)
		wantChanged bool
	}{
		{
			name: "edited manifest",
			edit: func(t *testing.T, location string) {
				path := filepath.Join(location, "units", "web", "kubernetes", "service.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				write(t, path, strings.Replace(string(data), "8080", "8081", 1))
			},
			wantChanged: true,
		},
		{
			name: "edited Terraform root",
			edit: func(t *testing.T, location string) {
				path := filepath.Join(location, "terraform", "main.tf.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				write(t, path, strings.Replace(string(data), `"10"`, `"11"`, 1))
			},
			wantChanged: true,
		},
		{
			name: "additional YAML manifest",
			edit: func(t *testing.T, location string) {
				write(t, filepath.Join(location, "units", "web", "kubernetes", "extra.yaml"), "apiVersion: v1\nkind: Service\nmetadata:\n  name: extra\n")
			},
			wantChanged: true,
		},
		{
			name: "notes beside the manifests",
			edit: func(t *testing.T, location string) {
				write(t, filepath.Join(location, "units", "web", "kubernetes", "notes.txt"), "not a manifest")
			},
		},
		{
			name: "edited README",
			edit: func(t *testing.T, location string) {
				write(t, filepath.Join(location, "README.md"), "# edited\n")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner, location := generatedRuntime(t)

			_, before, err := loadAndSnapshot(t, owner)
			if err != nil {
				t.Fatal(err)
			}

			tt.edit(t, location)

			_, after, err := loadAndSnapshot(t, owner)
			if err != nil {
				t.Fatal(err)
			}

			changed := after.Hash != before.Hash
			compatibleChanged := after.CompatibleHashes[0] != before.CompatibleHashes[0]

			if changed != tt.wantChanged || compatibleChanged != tt.wantChanged {
				t.Fatalf("fingerprint changed %t (compatible %t), want %t", changed, compatibleChanged, tt.wantChanged)
			}
		})
	}
}

func TestFiles_SnapshotReadsYAMLStreams(t *testing.T) {
	t.Parallel()

	owner, location := generatedRuntime(t)
	write(t, filepath.Join(location, "units", "web", "kubernetes", "more.yml"), `---
apiVersion: v1
kind: Service
metadata:
  name: admin
---
# an empty document is skipped
---
apiVersion: v1
kind: Service
metadata:
  name: metrics
spec:
  ports:
    - port: 9090
`)

	_, snapshot, err := loadAndSnapshot(t, owner)
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	for _, manifest := range snapshot.Manifests {
		names = append(names, manifest.Key())
	}

	if !slices.Contains(names, "Service/admin") || !slices.Contains(names, "Service/metrics") || len(names) != 5 {
		t.Fatalf("manifests = %v, want both YAML documents and no empty one", names)
	}
}

func TestFiles_SnapshotRefuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, owner project.Project, location string)
	}{
		{
			name: "missing component",
			arrange: func(t *testing.T, owner project.Project, _ string) {
				if err := os.RemoveAll(filepath.Join(owner.Root, "services", "web")); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "component that is a file",
			arrange: func(t *testing.T, owner project.Project, _ string) {
				if err := os.RemoveAll(filepath.Join(owner.Root, "jobs", "batch")); err != nil {
					t.Fatal(err)
				}
				write(t, filepath.Join(owner.Root, "jobs", "batch"), "")
			},
		},
		{
			name: "symlinked manifest",
			arrange: func(t *testing.T, _ project.Project, location string) {
				if err := os.Symlink(filepath.Join(location, "units", "web", "kubernetes", "service.json"),
					filepath.Join(location, "units", "web", "kubernetes", "linked.json")); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "manifest above the size bound",
			arrange: func(t *testing.T, _ project.Project, location string) {
				write(t, filepath.Join(location, "units", "web", "kubernetes", "huge.json"), `{"padding": "`+strings.Repeat("x", 1<<20)+`"}`)
			},
		},
		{
			name: "undecodable manifest",
			arrange: func(t *testing.T, _ project.Project, location string) {
				write(t, filepath.Join(location, "units", "web", "kubernetes", "broken.yaml"), "kind: [unterminated")
			},
		},
		{
			name: "missing Terraform root",
			arrange: func(t *testing.T, _ project.Project, location string) {
				if err := os.RemoveAll(filepath.Join(location, "terraform")); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "Terraform root outside the policy",
			arrange: func(t *testing.T, _ project.Project, location string) {
				write(t, filepath.Join(location, "terraform", "extra.tf.json"), `{"module": {"x": {"source": "./x"}}}`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner, location := generatedRuntime(t)
			tt.arrange(t, owner, location)

			if _, _, err := loadAndSnapshot(t, owner); err == nil {
				t.Fatal("Snapshot() accepted the configuration")
			}
		})
	}
}

func TestFiles_SnapshotChecksDeclaredHooks(t *testing.T) {
	t.Parallel()

	owner, location := generatedRuntime(t)

	descriptor, err := Files{}.Load(owner, identifier(t, "demo"))
	if err != nil {
		t.Fatal(err)
	}

	descriptor.Hooks = map[microruntime.Identifier]microruntime.RelativePath{identifier(t, "seed"): relativePath(t, "hooks/seed.sh")}

	if _, err := (Files{}).Snapshot(owner, descriptor); err == nil {
		t.Fatal("Snapshot() accepted a missing hook")
	}

	write(t, filepath.Join(location, "hooks", "seed.sh"), "#!/bin/sh\n")

	if _, err := (Files{}).Snapshot(owner, descriptor); err != nil {
		t.Fatalf("Snapshot() with the hook present = %v", err)
	}
}
