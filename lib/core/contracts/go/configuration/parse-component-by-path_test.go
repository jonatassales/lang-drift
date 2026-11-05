package configuration

import "testing"

func TestParse_ComponentByPath(t *testing.T) {
	config, err := Parse(validConfig().Markdown())
	if err != nil {
		t.Fatal(err)
	}

	component := config.Components[0]
	defaulted := component.DockerfileName() == "Dockerfile" &&
		component.Port == 0 &&
		component.Readiness == nil

	if !defaulted {
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
