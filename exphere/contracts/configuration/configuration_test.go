package configuration

import (
	"errors"
	"strings"
	"testing"
)

func validConfig() Configuration {
	return Configuration{Version: 1, RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "project", Revision: 1, Name: "Example", CPUMillis: 1000, MemoryMiB: 512, Model: "test-model", Effort: "medium", Components: []Component{{Path: ".", Image: "api:local", Port: 8080, ReadinessPath: "/health", Build: &Build{Executable: "docker", Args: []string{"build", "-t", "api:local", "."}}}}}
}
func TestParse_ExplicitMarkdownContract(t *testing.T) {
	original := validConfig().Markdown()
	for _, test := range []struct {
		name  string
		input string
		fail  bool
	}{{"root source", original, false}, {"unknown field", strings.Replace(original, `"version": 1`, `"extra": true, "version": 1`, 1), true}, {"ambiguous fences", original + original, true}, {"prose alone", "Run this shell command", true}, {"escaping source", strings.Replace(original, `"path": "."`, `"path": "../api"`, 1), true}, {"zero CPU", strings.Replace(original, `"cpuMillis": 1000`, `"cpuMillis": 0`, 1), true}, {"unbounded CPU", strings.Replace(original, `"cpuMillis": 1000`, `"cpuMillis": 64001`, 1), true}} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.input)
			if (err != nil) != test.fail {
				t.Fatalf("Parse err=%v wantFail=%v", err, test.fail)
			}
		})
	}
}
func TestConfiguration_DuplicateBasenamesRemainSeparate(t *testing.T) {
	config := validConfig()
	component := config.Components[0]
	component.Path = "apps/api"
	config.Components[0].Path = "services/api"
	config.Components = append(config.Components, component)
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	config.Components[1].Path = "services/api"
	if config.Validate() == nil {
		t.Fatal("duplicate source accepted")
	}
}

func TestBuild_RecipePolicy(t *testing.T) {
	refused := []Build{
		{Executable: "python3", Args: []string{"-c", "print(1)"}},
		{Executable: "node", Args: []string{"-e", "fetch('x')"}},
		{Executable: "sh", Args: []string{"-c", "id"}},
		{Executable: "sh", Args: []string{"build.sh"}},
		{Executable: "bash", Args: []string{"build.sh"}},
		{Executable: "curl", Args: []string{"example.com"}},
		{Executable: "npx", Args: []string{"vite", "build"}},
		{Executable: "docker", Args: []string{"run", "api:local"}},
		{Executable: "go", Args: []string{"-modfile", "x"}},
		{Executable: "make", Args: []string{"-j4"}},
		{Executable: "build.sh"},
		{Executable: "/usr/local/bin/build.sh"},
		{Executable: "../build.sh"},
		{Executable: "./tools/../../build.sh"},
		{Executable: "./build.sh", Args: []string{"--output", "/tmp"}},
		{Executable: `.\build.cmd`},
	}
	for _, build := range refused {
		if err := build.Validate(); !errorsIsRefused(err) {
			t.Fatalf("%v accepted, want ErrRecipeRefused", build)
		}
	}

	accepted := []Build{
		{Executable: "./build.sh"},
		{Executable: "./scripts/release.sh", Args: []string{"--prod"}},
		{Executable: "docker", Args: []string{"build", "-t", "api:local", "."}},
		{Executable: "go", Args: []string{"build", "./..."}},
		{Executable: "npm", Args: []string{"run", "build"}},
		{Executable: "pnpm", Args: []string{"build"}},
		{Executable: "make", Args: []string{"image"}},
		{Executable: "make"},
	}
	for _, build := range accepted {
		if err := build.Validate(); err != nil {
			t.Fatalf("%v refused: %v", build, err)
		}
	}
}

func errorsIsRefused(err error) bool { return errors.Is(err, ErrRecipeRefused) }
