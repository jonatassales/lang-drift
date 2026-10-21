package connection

import (
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
