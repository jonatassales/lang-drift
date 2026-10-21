package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
)

func TestStartRuntime_StartsAndRecordsTheRuntime(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	outcome := w.up(t)

	if outcome.Status != "running" || !outcome.Settled || outcome.Operation != "up" {
		t.Fatalf("outcome = %s, settled %t, operation %s; want a settled running up", outcome.Status, outcome.Settled, outcome.Operation)
	}

	if outcome.Namespace != namespaceName || outcome.NamespaceUID == "" || outcome.ClusterUID != kubeSystemUID ||
		outcome.Context != w.clusters.record.Context || outcome.ConfigHash != "hash-current" {
		t.Errorf("outcome target = %+v, want the recorded cluster, namespace, and configuration", outcome)
	}

	if len(outcome.Resources) != 3 {
		t.Fatalf("outcome records %d resources, want 3", len(outcome.Resources))
	}

	for _, resource := range outcome.Resources {
		if resource.UID == "" || resource.Unit == "" || resource.Namespace != namespaceName {
			t.Errorf("resource %+v lacks its created identity", resource)
		}
	}

	if outcome.Units["web"] != "completed" || outcome.Units["batch"] != "completed" {
		t.Errorf("units = %v, want every unit completed", outcome.Units)
	}

	for _, image := range outcome.Images {
		if image.ImageStatus != microruntime.ImageObserved || image.PodUID == "" {
			t.Errorf("image %+v was not observed on a live pod", image)
		}
	}

	if len(outcome.Images) != 2 {
		t.Errorf("observed %d containers, want web and batch", len(outcome.Images))
	}

	if outcome.ToolVersions["kubectl"].Status != "observed" || outcome.ToolVersions["terraform"].Status != "unavailable" {
		t.Errorf("tool versions = %v, want observed and unavailable tools distinguished", outcome.ToolVersions)
	}

	record := w.ledger.last(t, "demo")
	encoded, err := json.Marshal(outcome)
	if err != nil {
		t.Fatal(err)
	}

	if len(record.Token) != 32 || strings.Contains(string(encoded), record.Token) {
		t.Error("the private ownership token is missing from the ledger or exposed in the outcome")
	}

	first := w.ledger.saved[0]
	if first.Status != microruntime.StatusPreparing || first.NamespaceUID != "" || len(first.Resources) != 3 || first.Settled {
		t.Errorf("first saved record = %+v, want every obligation recorded before preparation", first)
	}

	if len(w.preparation.applied) != 1 {
		t.Fatalf("preparation ran %d times, want once", len(w.preparation.applied))
	}

	prepared := w.preparation.applied[0]
	if prepared.Token != record.Token || prepared.Namespace != namespaceName || prepared.Owner != "team-payments" ||
		prepared.Access.Kubeconfig != kubeconfigPath || prepared.Access.Context != w.clusters.record.Context ||
		prepared.Runtime.String() != "demo" || !strings.HasSuffix(prepared.Directory, "/terraform") {
		t.Errorf("preparation request = %+v, want the runtime's explicit target and token", prepared)
	}

	if !slices.Equal(w.kubernetes.dryRuns, []string{"Deployment/web", "Service/web", "Job/batch"}) ||
		!slices.Equal(w.kubernetes.creates, w.kubernetes.dryRuns) {
		t.Errorf("dry runs %v and creates %v, want every object validated, then created in order", w.kubernetes.dryRuns, w.kubernetes.creates)
	}

	if len(w.images.checked) != 1 || !slices.Equal(w.images.checked[0], []string{"batch:dev", "web:dev"}) {
		t.Errorf("images checked = %v, want both prebuilt images", w.images.checked)
	}

	if len(w.evidence.executions) != 1 || w.evidence.executions[0].RunID != outcome.RunID || w.evidence.executions[0].Status != "running" {
		t.Errorf("evidence = %+v, want the final outcome recorded once", w.evidence.executions)
	}

	w.requireLockBalanced(t)
}

func TestStartRuntime_RefusesBeforeRecording(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		runtime string
		arrange func(*testing.T, *world)
		wantErr error
	}{
		{name: "invalid runtime ID", runtime: "Demo", arrange: func(*testing.T, *world) {}, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unknown project", arrange: func(_ *testing.T, w *world) { w.projects.err = application.NotFound("no project") }, wantErr: application.ErrNotFound},
		{name: "operation already in progress", arrange: func(_ *testing.T, w *world) { w.locks.err = errInjected }, wantErr: errInjected},
		{
			name:    "descriptor for another runtime",
			arrange: func(t *testing.T, w *world) { w.configuration.descriptor.ID = mustIdentifier(t, "other") },
			wantErr: microruntime.ErrInvalidDescriptor,
		},
		{
			name: "manifest outside the policy",
			arrange: func(t *testing.T, w *world) {
				w.configuration.manifests = func() []microruntime.Manifest {
					manifests := demoManifests(t)
					manifests[2].Object["kind"] = "CronJob"

					return manifests
				}
			},
			wantErr: microruntime.ErrUnsupportedManifest,
		},
		{name: "no cluster", arrange: func(_ *testing.T, w *world) { w.clusters.record = nil }, wantErr: application.ErrNotFound},
		{name: "cluster not ready", arrange: func(_ *testing.T, w *world) { w.clusters.record.Status = cluster.StatusFailed }, wantErr: cluster.ErrNotReady},
		{
			name:    "replaced cluster",
			arrange: func(_ *testing.T, w *world) { w.observer.namespaces["kube-system"] = "uid-replacement" },
			wantErr: cluster.ErrUntrustedIdentity,
		},
		{
			name:    "context of another cluster",
			arrange: func(_ *testing.T, w *world) { w.configuration.descriptor.Target.Context = "kind-other" },
			wantErr: application.ErrTargetMismatch,
		},
		{name: "images not loaded", arrange: func(_ *testing.T, w *world) { w.images.missing = errInjected }, wantErr: errInjected},
		{name: "runtime still running", arrange: func(t *testing.T, w *world) { w.up(t) }, wantErr: microruntime.ErrExistingObligations},
		{
			name: "previous record is inconsistent",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.ledger.last(t, "demo").Token = ""
			},
			wantErr: microruntime.ErrInconsistentExecution,
		},
		{
			name: "namespace that EX did not create",
			arrange: func(_ *testing.T, w *world) {
				w.kubernetes.put(microruntime.Object{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": namespaceName, "uid": "uid-foreign"}})
			},
			wantErr: application.ErrAdoptionRefused,
		},
		{name: "ledger cannot record obligations", arrange: func(_ *testing.T, w *world) { w.ledger.saveErr = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)

			applied, creates, evidence := len(w.preparation.applied), len(w.kubernetes.creates), len(w.evidence.executions)

			request := w.runtimeRequest()
			if tt.runtime != "" {
				request.Runtime = tt.runtime
			}

			outcome, err := w.startRuntime().Execute(context.Background(), request)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if outcome.RunID != "" || outcome.Status != "" {
				t.Errorf("Execute() returned outcome %+v with an error, want the zero value", outcome)
			}

			if len(w.preparation.applied) != applied || len(w.kubernetes.creates) != creates || len(w.evidence.executions) != evidence {
				t.Error("a refused up prepared, created, or recorded something")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestStartRuntime_SettlesEveryFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		arrange    func(*world)
		timeout    time.Duration
		wantErr    error
		wantStatus microruntime.ExecutionStatus
		wantUnits  map[string]microruntime.UnitStatus
		wantCreate []string
		wantImages bool
	}{
		{
			name:       "preparation fails",
			arrange:    func(w *world) { w.preparation.applyErr = errInjected },
			wantErr:    errInjected,
			wantStatus: microruntime.StatusFailed,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitNotAttempted, "batch": microruntime.UnitNotAttempted},
		},
		{
			name:       "dry run rejects an object",
			arrange:    func(w *world) { w.kubernetes.dryRunErr = errInjected },
			wantErr:    errInjected,
			wantStatus: microruntime.StatusFailed,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitFailed, "batch": microruntime.UnitNotAttempted},
		},
		{
			name: "object already exists in the new namespace",
			arrange: func(w *world) {
				w.preparation.extra = []microruntime.Object{{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"name": "web"}}}
			},
			wantErr:    application.ErrAdoptionRefused,
			wantStatus: microruntime.StatusFailed,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitNotAttempted, "batch": microruntime.UnitNotAttempted},
		},
		{
			name:       "API server rejects a creation",
			arrange:    func(w *world) { w.kubernetes.createErr = errInjected },
			wantErr:    errInjected,
			wantStatus: microruntime.StatusFailed,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitFailed, "batch": microruntime.UnitNotAttempted},
			wantCreate: []string{"Deployment/web"},
		},
		{
			name:       "Job fails",
			arrange:    func(w *world) { w.kubernetes.jobFails = true },
			wantErr:    microruntime.ErrWorkloadFailed,
			wantStatus: microruntime.StatusFailed,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitCompleted, "batch": microruntime.UnitFailed},
			wantCreate: []string{"Deployment/web", "Service/web", "Job/batch"},
			wantImages: true,
		},
		{
			name:       "deadline passes while waiting for readiness",
			arrange:    func(w *world) { w.kubernetes.unready = true },
			timeout:    50 * time.Millisecond,
			wantErr:    context.DeadlineExceeded,
			wantStatus: microruntime.StatusInterrupted,
			wantUnits:  map[string]microruntime.UnitStatus{"web": microruntime.UnitFailed, "batch": microruntime.UnitFailed},
			wantCreate: []string{"Deployment/web", "Service/web", "Job/batch"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)

			ctx := context.Background()
			if tt.timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}

			_, err := w.startRuntime().Execute(ctx, w.runtimeRequest())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			record := w.ledger.last(t, "demo")
			if record.Status != tt.wantStatus || !record.Settled {
				t.Fatalf("record = %s, settled %t; want %s and settled", record.Status, record.Settled, tt.wantStatus)
			}

			for unit, want := range tt.wantUnits {
				if got := record.Units[mustIdentifier(t, unit)]; got != want {
					t.Errorf("unit %s = %s, want %s", unit, got, want)
				}
			}

			if !slices.Equal(w.kubernetes.creates, tt.wantCreate) {
				t.Errorf("creates = %v, want %v", w.kubernetes.creates, tt.wantCreate)
			}

			observed := slices.ContainsFunc(record.Images, func(image microruntime.ContainerObservation) bool { return image.ImageID != "" })
			if observed != tt.wantImages {
				t.Errorf("observed images recorded = %t, want %t", observed, tt.wantImages)
			}

			if len(w.evidence.executions) != 1 || w.evidence.executions[0].Status != string(tt.wantStatus) {
				t.Errorf("evidence = %+v, want the settled outcome", w.evidence.executions)
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestStartRuntime_ReportsUnrecordedEvidence(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.evidence.err = errInjected

	outcome, err := w.startRuntime().Execute(context.Background(), w.runtimeRequest())
	if !errors.Is(err, errInjected) || outcome.RunID != "" {
		t.Fatalf("Execute() = %+v, %v; want the evidence failure and no outcome", outcome, err)
	}

	if record := w.ledger.last(t, "demo"); record.Status != microruntime.StatusRunning || !record.Settled {
		t.Fatalf("record = %s, settled %t; the ledger must still settle", record.Status, record.Settled)
	}
}

func TestStartRuntime_StartsAgainAfterCleanup(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	first := w.up(t)

	if _, err := w.stopRuntime().Execute(context.Background(), application.StopRuntimeRequest{RuntimeRequest: w.runtimeRequest()}); err != nil {
		t.Fatalf("runtime down: %v", err)
	}

	firstToken := w.ledger.last(t, "demo").Token
	second := w.up(t)

	if second.RunID == first.RunID || w.ledger.last(t, "demo").Token == firstToken {
		t.Fatal("a new up reused the previous run ID or ownership token")
	}
}

func TestStartRuntime_ReportsALedgerThatStopsRecording(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.ledger.failAfter = 1

	if _, err := w.startRuntime().Execute(context.Background(), w.runtimeRequest()); !errors.Is(err, errInjected) {
		t.Fatalf("Execute() error = %v, want the ledger failure", err)
	}

	if record := w.ledger.last(t, "demo"); record.Status != microruntime.StatusPreparing || record.Settled {
		t.Errorf("record = %s, settled %t; want the obligations recorded before the ledger failed", record.Status, record.Settled)
	}

	if len(w.evidence.executions) != 1 || w.evidence.executions[0].Status != "failed" {
		t.Errorf("evidence = %+v, want the failure recorded even though the ledger is unavailable", w.evidence.executions)
	}
}
