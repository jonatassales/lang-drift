package configuration

import (
	"errors"
	"strings"
	"testing"
)

func validConfig() Configuration {
	return Configuration{
		Version:    2,
		RuntimeUID: "runtime-uid",
		RuntimeID:  "runtime",
		ProjectID:  "project",
		Revision:   1,
		Name:       "Example",
		CPUMillis:  1000,
		MemoryMiB:  512,
		Model:      "test-model",
		Effort:     "medium",
		Components: []Component{
			{
				Name: "api",
				Path: ".",
			},
		},
	}
}

func TestParse_ExplicitMarkdownContract(t *testing.T) {
	original := validConfig().Markdown()

	for _, test := range []struct {
		name   string
		input  string
		reason string
	}{
		{"root source", original, ""},
		{"unknown field", strings.Replace(original, `"version": 2`, `"extra": true, "version": 2`, 1), "re-plan the runtime"},
		{"version 1", strings.Replace(original, `"version": 2`, `"version": 1`, 1), "re-plan the runtime"},
		{"version 3", strings.Replace(original, `"version": 2`, `"version": 3`, 1), "re-plan the runtime"},
		{"inline image", strings.Replace(original, `"path": "."`, `"path": ".", "image": "api:local"`, 1), "re-plan the runtime"},
		{"build recipe", strings.Replace(original, `"path": "."`, `"path": ".", "build": {"executable": "make", "args": []}`, 1), "re-plan the runtime"},
		{"readiness path", strings.Replace(original, `"path": "."`, `"path": ".", "readinessPath": "/health"`, 1), "re-plan the runtime"},
		{"ambiguous fences", original + original, "exactly one langdrift JSON fence"},
		{"prose alone", "Run this shell command", "exactly one langdrift JSON fence"},
		{"escaping source", strings.Replace(original, `"path": "."`, `"path": "../api"`, 1), "clean relative path"},
		{"zero CPU", strings.Replace(original, `"cpuMillis": 1000`, `"cpuMillis": 0`, 1), "CPU and memory"},
		{"unbounded CPU", strings.Replace(original, `"cpuMillis": 1000`, `"cpuMillis": 64001`, 1), "CPU and memory"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.input)

			if test.reason == "" && err != nil {
				t.Fatalf("Parse refused a valid fence: %v", err)
			}

			if test.reason != "" && (err == nil || !strings.Contains(err.Error(), test.reason)) {
				t.Fatalf("Parse err = %v, want %q", err, test.reason)
			}
		})
	}
}

func TestParse_ComponentByPath(t *testing.T) {
	config, err := Parse(validConfig().Markdown())
	if err != nil {
		t.Fatal(err)
	}

	component := config.Components[0]
	if component.DockerfileName() != "Dockerfile" || component.Port != 0 || component.Readiness != nil {
		t.Fatalf("component = %+v, want the Dockerfile default and no port or readiness", component)
	}

	if component.DockerfilePath() != "Dockerfile" {
		t.Fatalf("DockerfilePath = %q, want the root Dockerfile", component.DockerfilePath())
	}

	component.Path = "apps/api"
	component.Dockerfile = "Dockerfile.local"

	if component.DockerfilePath() != "apps/api/Dockerfile.local" {
		t.Fatalf("DockerfilePath = %q, want the named file in the folder", component.DockerfilePath())
	}
}

func TestConfiguration_ComponentRules(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*Component)
		reason string
	}{
		{"every refinement", func(component *Component) {
			component.Dockerfile = "Dockerfile.dev"
			component.Port = 8080
			component.Readiness = &Readiness{
				Kind: ReadinessHTTP,
				Path: "/health",
			}
			component.EnvFile = "apps/api/runtime.env"
			component.PassEnv = []string{"API_TOKEN"}
		}, ""},
		{"tcp readiness", func(component *Component) {
			component.Readiness = &Readiness{Kind: ReadinessTCP}
		}, ""},
		{"no name", func(component *Component) {
			component.Name = ""
		}, "DNS-label"},
		{"name not a DNS label", func(component *Component) {
			component.Name = "Api_Server"
		}, "DNS-label"},
		{"dockerfile with a slash", func(component *Component) {
			component.Dockerfile = "build/Dockerfile"
		}, "file name inside the component folder"},
		{"dockerfile parent", func(component *Component) {
			component.Dockerfile = ".."
		}, "file name inside the component folder"},
		{"port out of range", func(component *Component) {
			component.Port = 70000
		}, "port"},
		{"unknown readiness", func(component *Component) {
			component.Readiness = &Readiness{Kind: "grpc"}
		}, "readiness kind is tcp or http"},
		{"http without a path", func(component *Component) {
			component.Readiness = &Readiness{Kind: ReadinessHTTP}
		}, "path starting with /"},
		{"tcp with a path", func(component *Component) {
			component.Readiness = &Readiness{
				Kind: ReadinessTCP,
				Path: "/health",
			}
		}, "tcp readiness takes no path"},
		{"env file outside the checkout", func(component *Component) {
			component.EnvFile = "../secrets.env"
		}, "clean relative path"},
		{"env file in .langdrift", func(component *Component) {
			component.EnvFile = ".langdrift/runtime.env"
		}, "outside .langdrift"},
		{"lowercase pass-through", func(component *Component) {
			component.PassEnv = []string{"api_token"}
		}, "uppercase variable name"},
		{"repeated pass-through", func(component *Component) {
			component.PassEnv = []string{"API_TOKEN", "API_TOKEN"}
		}, "named twice"},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := validConfig()
			test.change(&config.Components[0])

			err := config.Validate()

			if test.reason == "" && err != nil {
				t.Fatalf("Validate refused: %v", err)
			}

			if test.reason != "" && (err == nil || !strings.Contains(err.Error(), test.reason)) {
				t.Fatalf("Validate err = %v, want %q", err, test.reason)
			}

			if err != nil && !errors.Is(err, ErrInvalidConfiguration) {
				t.Fatalf("Validate err = %v, want ErrInvalidConfiguration", err)
			}
		})
	}
}

func TestConfiguration_NamesAndPathsAreUnique(t *testing.T) {
	config := validConfig()
	config.Components[0].Path = "services/api"
	config.Components = append(config.Components, Component{
		Name: "web",
		Path: "apps/api",
	})

	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}

	config.Components[1].Path = "services/api"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "unique source paths") {
		t.Fatalf("duplicate source accepted: %v", err)
	}

	config.Components[1].Path = "apps/api"
	config.Components[1].Name = "api"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "unique names") {
		t.Fatalf("duplicate name accepted: %v", err)
	}
}

func TestComponent_CheckEnv(t *testing.T) {
	component := Component{
		Name:    "api",
		Path:    "apps/api",
		EnvFile: "apps/api/runtime.env",
		PassEnv: []string{"API_TOKEN"},
	}

	if err := component.CheckEnv(map[string]string{"LOG_URL": "http://log:3500"}); err != nil {
		t.Fatalf("CheckEnv refused distinct names: %v", err)
	}

	err := component.CheckEnv(map[string]string{"API_TOKEN": "committed"})
	if err == nil || !strings.Contains(err.Error(), "API_TOKEN") {
		t.Fatalf("CheckEnv err = %v, want the overlapping name", err)
	}
}
