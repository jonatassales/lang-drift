package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

type fakeAgents struct {
	err        error
	scaffolded []string
}

func (agents *fakeAgents) Scaffold(_ project.Project, name microruntime.Identifier) (string, error) {
	agents.scaffolded = append(agents.scaffolded, name.String())

	return projectRoot + "/.exphere/agents/" + name.String(), agents.err
}

func TestInitializeAgent_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		agent     string
		projects  error
		scaffold  error
		wantErr   error
		wantAgent bool
	}{
		{name: "new agent", agent: "researcher", wantAgent: true},
		{name: "invalid name", agent: "Researcher", wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unknown project", agent: "researcher", projects: errInjected, wantErr: errInjected},
		{name: "existing agent", agent: "researcher", scaffold: errInjected, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.projects.err = tt.projects
			agents := &fakeAgents{err: tt.scaffold}

			result, err := application.InitializeAgent{Projects: w.projects, Agents: agents}.Execute(
				context.Background(), application.InitializeAgentRequest{Project: w.runtimeRequest().Project, Name: tt.agent})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantAgent && (result.Agent != "researcher" || result.Path != projectRoot+"/.exphere/agents/researcher") {
				t.Fatalf("result = %+v, want the published agent", result)
			}
		})
	}
}
