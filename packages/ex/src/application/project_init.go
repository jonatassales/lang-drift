package application

import (
	"context"
	"errors"

	"github.com/jonatassales/tao/exphere/cli/src/domain/project"
)

// InitializeProject records a project's identity in a directory, or reports
// the project already there. It never installs tools or starts infrastructure.
type InitializeProject struct {
	Projects  ProjectOpener
	Creator   ProjectCreator
	Bootstrap ProjectBootstrapper

	// Injectable for tests; nil uses the real sources.
	IDs IDSource
}

// InitializeProjectRequest names the directory that becomes the project root.
type InitializeProjectRequest struct {
	Directory string
}

// ProjectResult describes an initialized project.
type ProjectResult struct {
	Root   string `json:"project"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Execute reuses existing metadata, or creates it with a fresh identity.
func (useCase InitializeProject) Execute(_ context.Context, request InitializeProjectRequest) (ProjectResult, error) {
	selection := ProjectSelection{Directory: request.Directory, Explicit: true}

	existing, err := useCase.Projects.Open(selection)
	if err == nil {
		return useCase.prepared(existing, "already initialized")
	}

	if !errors.Is(err, ErrNotFound) {
		return ProjectResult{}, err
	}

	created, err := useCase.create(request.Directory)

	switch {
	case errors.Is(err, ErrAlreadyExists):
		// A concurrent initialization won the race; report its project.
		existing, err := useCase.Projects.Open(selection)
		if err != nil {
			return ProjectResult{}, err
		}

		return useCase.prepared(existing, "already initialized")

	case err != nil:
		return ProjectResult{}, err
	}

	return useCase.prepared(created, "initialized")
}

func (useCase InitializeProject) create(directory string) (project.Project, error) {
	text, err := useCase.IDs.id()
	if err != nil {
		return project.Project{}, err
	}

	id, err := project.ParseID(text)
	if err != nil {
		return project.Project{}, err
	}

	return useCase.Creator.Create(directory, id)
}

func projectResultOf(owner project.Project, status string) ProjectResult {
	return ProjectResult{Root: owner.Root, ID: owner.ID.String(), Status: status}
}

// ProjectBootstrapper adds the connected agent and reference workflow without
// replacing existing consumer-authored files or running dependencies.
type ProjectBootstrapper interface{ Prepare(project.Project) error }

func (useCase InitializeProject) prepared(owner project.Project, status string) (ProjectResult, error) {
	if useCase.Bootstrap != nil {
		if err := useCase.Bootstrap.Prepare(owner); err != nil {
			return ProjectResult{}, err
		}
	}
	return projectResultOf(owner, status), nil
}
