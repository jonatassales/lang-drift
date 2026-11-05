package configuration

import (
	"strings"
	"testing"
)

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
