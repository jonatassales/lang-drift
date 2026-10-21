package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/project"
)

type fakeCreator struct {
	err     error
	created []project.Project
}

func (creator *fakeCreator) Create(directory string, id project.ID) (project.Project, error) {
	if creator.err != nil {
		return project.Project{}, creator.err
	}

	created := project.Project{ID: id, Root: directory}
	creator.created = append(creator.created, created)

	return created, nil
}

// sequenceProjects answers each Open with the next prepared result.
type sequenceProjects struct {
	results []error
	project project.Project
	opened  []application.ProjectSelection
}

func (projects *sequenceProjects) Open(selection application.ProjectSelection) (project.Project, error) {
	projects.opened = append(projects.opened, selection)
	err := projects.results[0]
	projects.results = projects.results[1:]

	if err != nil {
		return project.Project{}, err
	}

	return projects.project, nil
}

func TestInitializeProject_Execute(t *testing.T) {
	t.Parallel()

	existing := newWorld(t).owner
	notFound := application.NotFound("no Exphere project at %s", projectRoot)

	tests := []struct {
		name        string
		opens       []error
		createErr   error
		wantStatus  string
		wantCreated bool
		wantErr     error
	}{
		{name: "new project", opens: []error{notFound}, wantStatus: "initialized", wantCreated: true},
		{name: "existing project", opens: []error{nil}, wantStatus: "already initialized"},
		{
			name:       "concurrent initialization won the race",
			opens:      []error{notFound, nil},
			createErr:  application.AlreadyExists("project metadata already exists at %s", projectRoot),
			wantStatus: "already initialized",
		},
		{name: "unreadable metadata", opens: []error{errInjected}, wantErr: errInjected},
		{name: "creation fails", opens: []error{notFound}, createErr: errInjected, wantErr: errInjected},
		{
			name:      "winner unreadable after the race",
			opens:     []error{notFound, errInjected},
			createErr: application.AlreadyExists("project metadata already exists"),
			wantErr:   errInjected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			projects := &sequenceProjects{results: tt.opens, project: existing}
			creator := &fakeCreator{err: tt.createErr}
			useCase := application.InitializeProject{Projects: projects, Creator: creator}

			result, err := useCase.Execute(context.Background(), application.InitializeProjectRequest{Directory: projectRoot})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if result.Status != tt.wantStatus || (len(creator.created) == 1) != tt.wantCreated {
				t.Fatalf("Execute() = %+v with %d creations, want %q", result, len(creator.created), tt.wantStatus)
			}

			for _, selection := range projects.opened {
				if !selection.Explicit || selection.Directory != projectRoot {
					t.Errorf("opened %+v, want the explicit directory only", selection)
				}
			}

			if tt.wantCreated {
				if _, err := project.ParseID(result.ID); err != nil || result.Root != projectRoot {
					t.Errorf("created %+v, want a fresh valid identity at the directory", result)
				}
			}
		})
	}
}
