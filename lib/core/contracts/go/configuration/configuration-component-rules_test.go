package configuration

import (
	"errors"
	"strings"
	"testing"
)

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

			missedReason := test.reason != "" &&
				(err == nil || !strings.Contains(err.Error(), test.reason))

			if missedReason {
				t.Fatalf("Validate err = %v, want %q", err, test.reason)
			}

			if err != nil && !errors.Is(err, ErrInvalidConfiguration) {
				t.Fatalf("Validate err = %v, want ErrInvalidConfiguration", err)
			}
		})
	}
}
