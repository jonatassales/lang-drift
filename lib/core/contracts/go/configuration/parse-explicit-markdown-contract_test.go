package configuration

import (
	"strings"
	"testing"
)

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

			missedReason := test.reason != "" &&
				(err == nil || !strings.Contains(err.Error(), test.reason))

			if missedReason {
				t.Fatalf("Parse err = %v, want %q", err, test.reason)
			}
		})
	}
}
