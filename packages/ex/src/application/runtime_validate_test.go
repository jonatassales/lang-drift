package application_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
)

func (w *world) validateRuntime() application.ValidateRuntime {
	return application.ValidateRuntime{Projects: w.projects, Configuration: w.configuration}
}

func TestValidateRuntime_SummarizesAValidConfiguration(t *testing.T) {
	t.Parallel()

	w := newWorld(t)

	result, err := w.validateRuntime().Execute(context.Background(), w.runtimeRequest())
	if err != nil {
		t.Fatal(err)
	}

	if result.Runtime != "demo" || !slices.Equal(result.Units, []string{"batch", "web"}) || result.Objects != 3 {
		t.Fatalf("result = %+v, want sorted units and every object counted", result)
	}

	if len(w.kubernetes.dryRuns) != 0 || len(w.preparation.applied) != 0 {
		t.Fatal("validation contacted the cluster or prepared something")
	}
}

func TestValidateRuntime_RejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		runtime string
		arrange func(*testing.T, *world)
		wantErr error
	}{
		{name: "invalid runtime ID", runtime: "demo/x", arrange: func(*testing.T, *world) {}, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unknown project", arrange: func(_ *testing.T, w *world) { w.projects.err = errInjected }, wantErr: errInjected},
		{
			name: "runtime not initialized",
			arrange: func(_ *testing.T, w *world) {
				w.configuration.loadErr = application.NotFound("runtime demo is not initialized")
			},
			wantErr: application.ErrNotFound,
		},
		{
			name:    "descriptor for another runtime",
			arrange: func(t *testing.T, w *world) { w.configuration.descriptor.ID = mustIdentifier(t, "other") },
			wantErr: microruntime.ErrInvalidDescriptor,
		},
		{
			name:    "incomplete descriptor",
			arrange: func(_ *testing.T, w *world) { w.configuration.descriptor.Owner = "" },
			wantErr: microruntime.ErrInvalidDescriptor,
		},
		{name: "unreadable native files", arrange: func(_ *testing.T, w *world) { w.configuration.snapshotErr = errInjected }, wantErr: errInjected},
		{
			name: "unit without objects",
			arrange: func(t *testing.T, w *world) {
				w.configuration.manifests = func() []microruntime.Manifest { return demoManifests(t)[:2] }
			},
			wantErr: microruntime.ErrUnsupportedManifest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)

			request := w.runtimeRequest()
			if tt.runtime != "" {
				request.Runtime = tt.runtime
			}

			if _, err := w.validateRuntime().Execute(context.Background(), request); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
