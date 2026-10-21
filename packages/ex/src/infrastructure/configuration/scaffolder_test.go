package configuration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// tempDir returns a fresh directory with symlinks resolved, because these
// adapters refuse symlinked ancestors such as macOS's /var link.
func tempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func identifier(t testing.TB, text string) microruntime.Identifier {
	t.Helper()

	parsed, err := microruntime.ParseIdentifier(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

func relativePath(t testing.TB, text string) microruntime.RelativePath {
	t.Helper()

	parsed, err := microruntime.ParseRelativePath(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

func quantity(t testing.TB, text string) microruntime.Quantity {
	t.Helper()

	parsed, err := microruntime.ParseQuantity(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

// testProject is a project whose component directories exist.
func testProject(t *testing.T) project.Project {
	t.Helper()

	id, err := project.ParseID("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}

	root := tempDir(t)
	for _, component := range []string{"services/web", "jobs/batch"} {
		if err := os.MkdirAll(filepath.Join(root, component), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	return project.Project{ID: id, Root: root}
}

func testBlueprint(t *testing.T) microruntime.Blueprint {
	t.Helper()

	resources := microruntime.Resources{
		CPURequest: quantity(t, "100m"), CPULimit: quantity(t, "500m"),
		MemoryRequest: quantity(t, "128Mi"), MemoryLimit: quantity(t, "256Mi"),
	}

	return microruntime.Blueprint{
		Owner:  "team-payments",
		Target: microruntime.KubernetesTarget{Context: "kind-ex-0123456789ab", Namespace: identifier(t, "demo-local")},
		Workloads: []microruntime.Workload{
			{
				Name: identifier(t, "web"), Component: relativePath(t, "services/web"), Image: "web:dev",
				Kind: microruntime.DeploymentWorkload, Resources: resources,
				Service: microruntime.ServiceSpec{Port: 8080, Replicas: 1, ReadinessPath: "/healthz"},
			},
			{
				Name: identifier(t, "batch"), Component: relativePath(t, "jobs/batch"), Image: "batch:dev",
				Kind: microruntime.JobWorkload, Resources: resources,
			},
		},
		Quotas: microruntime.ResourceQuotas{
			microruntime.QuotaRequestsCPU: quantity(t, "1"),
			microruntime.QuotaPods:        quantity(t, "10"),
		},
	}
}

// generatedRuntime publishes the guided configuration of runtime demo.
func generatedRuntime(t *testing.T) (project.Project, string) {
	t.Helper()

	owner := testProject(t)

	location, err := Scaffolder{}.Generate(owner, identifier(t, "demo"), testBlueprint(t))
	if err != nil {
		t.Fatal(err)
	}

	return owner, location
}

func relativeFiles(t *testing.T, root string) []string {
	t.Helper()

	var found []string

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		relative, err := filepath.Rel(root, path)
		found = append(found, filepath.ToSlash(relative))

		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	slices.Sort(found)

	return found
}

func TestScaffolder_Generate(t *testing.T) {
	t.Parallel()

	owner, location := generatedRuntime(t)

	if location != filepath.Join(owner.Root, ".exphere", "runtimes", "demo") {
		t.Fatalf("Generate() location = %q", location)
	}

	want := []string{
		"README.md",
		"schema.json",
		"terraform/main.tf.json",
		"units/batch/kubernetes/workload.json",
		"units/web/kubernetes/service.json",
		"units/web/kubernetes/workload.json",
	}
	if got := relativeFiles(t, location); !slices.Equal(got, want) {
		t.Fatalf("generated files = %v, want %v", got, want)
	}

	descriptor, err := Files{}.Load(owner, identifier(t, "demo"))
	if err != nil {
		t.Fatal(err)
	}

	if err := descriptor.Validate(); err != nil {
		t.Fatalf("the generated descriptor is invalid: %v", err)
	}

	snapshot, err := Files{}.Snapshot(owner, descriptor)
	if err != nil {
		t.Fatalf("the generated configuration cannot be snapshotted: %v", err)
	}

	if err := descriptor.ValidateManifests(snapshot.Manifests); err != nil {
		t.Fatalf("the generated manifests violate the policy: %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(location, "README.md"))
	if err != nil || !strings.HasPrefix(string(readme), "# demo\n") {
		t.Fatalf("README = %q, %v", readme, err)
	}
}

func TestScaffolder_GenerateRefuses(t *testing.T) {
	t.Parallel()

	t.Run("missing component directory", func(t *testing.T) {
		t.Parallel()

		owner := testProject(t)
		if err := os.RemoveAll(filepath.Join(owner.Root, "jobs")); err != nil {
			t.Fatal(err)
		}

		if _, err := (Scaffolder{}).Generate(owner, identifier(t, "demo"), testBlueprint(t)); err == nil {
			t.Fatal("Generate() referenced a missing component")
		}

		if _, err := os.Stat(filepath.Join(owner.Root, ".exphere", "runtimes", "demo")); !os.IsNotExist(err) {
			t.Fatalf("a refused generation wrote configuration: %v", err)
		}
	})

	t.Run("existing runtime", func(t *testing.T) {
		t.Parallel()

		owner, location := generatedRuntime(t)
		before := relativeFiles(t, location)

		if _, err := (Scaffolder{}).Generate(owner, identifier(t, "demo"), testBlueprint(t)); err == nil {
			t.Fatal("Generate() replaced an existing runtime")
		}

		if after := relativeFiles(t, location); !slices.Equal(before, after) {
			t.Fatalf("existing configuration changed: %v", after)
		}
	})

	t.Run("legacy runtime with the same ID", func(t *testing.T) {
		t.Parallel()

		owner := testProject(t)
		if err := os.MkdirAll(filepath.Join(owner.Root, ".micro-runtimes", "demo"), 0o755); err != nil {
			t.Fatal(err)
		}

		if _, err := (Scaffolder{}).Generate(owner, identifier(t, "demo"), testBlueprint(t)); err == nil || !strings.Contains(err.Error(), "legacy runtime") {
			t.Fatalf("Generate() error = %v, want the legacy runtime refusal", err)
		}
	})
}

func TestScaffolder_Scaffold(t *testing.T) {
	t.Parallel()

	owner := testProject(t)

	location, err := Scaffolder{}.Scaffold(owner, microruntime.Scaffold(identifier(t, "manual")))
	if err != nil {
		t.Fatal(err)
	}

	if got := relativeFiles(t, location); !slices.Equal(got, []string{"README.md", "schema.json"}) {
		t.Fatalf("scaffolded files = %v", got)
	}

	data, err := os.ReadFile(filepath.Join(location, "schema.json"))
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	if document["version"] != 2.0 || document["environment"] != "local" || document["owner"] != "" {
		t.Errorf("schema.json = %s, want an incomplete version 2 descriptor", data)
	}

	descriptor, err := Files{}.Load(owner, identifier(t, "manual"))
	if err != nil {
		t.Fatalf("the scaffold must decode so validation can explain what is missing: %v", err)
	}

	if err := descriptor.Validate(); err == nil {
		t.Fatal("the incomplete scaffold validated")
	}

	if _, err := (Scaffolder{}).Scaffold(owner, microruntime.Scaffold(identifier(t, "manual"))); err == nil {
		t.Fatal("Scaffold() replaced an existing runtime")
	}
}
