package application

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/domain/microruntime"
)

// InitializeAgent publishes an editable agent project inside the consumer
// project. It never installs dependencies, builds images, or runs the agent.
type InitializeAgent struct {
	Projects ProjectOpener
	Agents   AgentScaffolder
}

// InitializeAgentRequest names the agent to scaffold.
type InitializeAgentRequest struct {
	Project ProjectSelection
	Name    string
}

// AgentResult describes a scaffolded agent.
type AgentResult struct {
	Agent  string `json:"agent"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

// Execute validates the name and publishes the scaffold.
func (useCase InitializeAgent) Execute(_ context.Context, request InitializeAgentRequest) (AgentResult, error) {
	name, err := microruntime.ParseIdentifier(request.Name)
	if err != nil {
		return AgentResult{}, fmt.Errorf("agent name: %w", err)
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return AgentResult{}, err
	}

	location, err := useCase.Agents.Scaffold(owner, name)
	if err != nil {
		return AgentResult{}, err
	}

	return AgentResult{
		Agent:  name.String(),
		Path:   location,
		Status: "generated; building its image remains an explicit step",
	}, nil
}
