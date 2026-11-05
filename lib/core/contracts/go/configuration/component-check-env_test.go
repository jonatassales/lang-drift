package configuration

import (
	"strings"
	"testing"
)

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
