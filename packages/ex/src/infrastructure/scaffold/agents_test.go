package scaffold_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/scaffold"
)

func TestAgents_Scaffold(t *testing.T) {
	t.Parallel()

	owner := project.Project{Root: tempDir(t)}
	name, err := microruntime.ParseIdentifier("researcher")
	if err != nil {
		t.Fatal(err)
	}

	location, err := scaffold.Agents{}.Scaffold(owner, name)
	if err != nil {
		t.Fatal(err)
	}

	if location != filepath.Join(owner.Root, ".exphere", "agents", "researcher") {
		t.Fatalf("Scaffold() location = %q", location)
	}

	want := []string{
		".dockerignore", ".gitignore", "Dockerfile", "README.md", "package-lock.json",
		"package.json", "src/", "src/index.ts", "tsconfig.json",
	}
	if got := entries(t, location); !slices.Equal(got, want) {
		t.Fatalf("scaffolded entries = %v, want %v", got, want)
	}

	data, err := os.ReadFile(filepath.Join(location, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("package.json is not valid JSON: %v", err)
	}

	if _, err := (scaffold.Agents{}).Scaffold(owner, name); err == nil {
		t.Fatal("Scaffold() replaced an existing agent")
	}
}
