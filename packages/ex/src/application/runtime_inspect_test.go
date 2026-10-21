package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
)

func (w *world) inspectRuntime() application.InspectRuntime {
	return application.InspectRuntime{
		Projects: w.projects, Clusters: w.clusters, Observer: w.observer, Executions: w.ledger,
		Objects: w.kubernetes, Evidence: w.evidence, Tools: w.tools,
	}
}

func event(name, uid string) microruntime.Object {
	return microruntime.Object{
		"metadata":       map[string]any{"name": name},
		"reason":         "Scheduled",
		"type":           "Normal",
		"message":        "may echo SECRET=value",
		"count":          1.0,
		"involvedObject": map[string]any{"kind": "Pod", "name": name, "uid": uid},
	}
}

func TestInspectRuntime_ReportsTheRecordedRuntime(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)
	saves := len(w.ledger.saved)

	pod := w.kubernetes.live("Pod", "web-6f6-a")
	w.kubernetes.events = []microruntime.Object{event("web-6f6-a.1", pod.UID()), event("stranger.1", "uid-stranger")}
	w.kubernetes.metrics = []microruntime.Object{
		{"metadata": map[string]any{"name": "web-6f6-a"}, "timestamp": "now", "window": "30s",
			"containers": []any{map[string]any{"name": "web", "usage": map[string]any{"cpu": "1m"}}}},
		{"metadata": map[string]any{"name": "stranger"}},
	}

	report, err := w.inspectRuntime().Execute(context.Background(), w.runtimeRequest())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if report.Runtime != "demo" || report.ObservationID == "" || report.ExecutionRunID != w.ledger.last(t, "demo").RunID {
		t.Errorf("report identity = %+v", report)
	}

	if len(report.Resources) != 3 {
		t.Fatalf("report has %d resources, want 3", len(report.Resources))
	}

	for _, resource := range report.Resources {
		if resource.Observed != "present" || resource.ObservedUID != resource.Resource.UID {
			t.Errorf("resource %+v, want the recorded object present", resource)
		}
	}

	if len(report.Containers) != 2 || len(report.Quotas) != 1 || report.Quotas[0].Name != "runtime" {
		t.Errorf("containers %d, quotas %+v; want both containers and the runtime quota", len(report.Containers), report.Quotas)
	}

	if len(report.Events) != 1 || report.Events[0].Object.UID != pod.UID() {
		t.Errorf("events = %+v, want only the event about an owned object", report.Events)
	}

	if report.Usage.Status != "observed" || len(report.Usage.Pods) != 1 || report.Usage.Pods[0].Pod != "web-6f6-a" {
		t.Errorf("usage = %+v, want only the owned pod's metrics", report.Usage)
	}

	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}

	for _, secret := range []string{w.ledger.last(t, "demo").Token, "SECRET=value"} {
		if strings.Contains(string(encoded), secret) {
			t.Errorf("the report exposes %q", secret)
		}
	}

	if len(w.evidence.inspections) != 1 || w.evidence.inspections[0].ObservationID != report.ObservationID {
		t.Errorf("inspection evidence = %d reports, want this one", len(w.evidence.inspections))
	}

	if w.locks.locked != 1 || len(w.ledger.saved) != saves {
		t.Errorf("inspection took the project lock or changed the ledger")
	}
}

func TestInspectRuntime_ReportsAbsentObjectsAndUnavailableMetrics(t *testing.T) {
	t.Parallel()

	t.Run("namespace removed", func(t *testing.T) {
		t.Parallel()

		w := newWorld(t)
		w.up(t)
		w.kubernetes.objects = nil

		report, err := w.inspectRuntime().Execute(context.Background(), w.runtimeRequest())
		if err != nil {
			t.Fatal(err)
		}

		for _, resource := range report.Resources {
			if resource.Observed != "absent" {
				t.Errorf("resource %+v, want absent", resource)
			}
		}

		encoded, _ := json.Marshal(report.Resources[0])
		if strings.Contains(string(encoded), "observedUID") {
			t.Errorf("an absent resource serializes details: %s", encoded)
		}
	})

	t.Run("one object removed and no metrics API", func(t *testing.T) {
		t.Parallel()

		w := newWorld(t)
		w.up(t)
		w.kubernetes.metricsErr = errInjected
		w.kubernetes.edit("Job", "batch", func(object microruntime.Object) { object["kind"] = "Removed" })

		report, err := w.inspectRuntime().Execute(context.Background(), w.runtimeRequest())
		if err != nil {
			t.Fatal(err)
		}

		if report.Resources[2].Observed != "absent" || report.Resources[0].Observed != "present" {
			t.Errorf("resources = %+v, want the Job absent and the others present", report.Resources)
		}

		if report.Usage.Status != "unavailable" || report.Usage.Reason == "" {
			t.Errorf("usage = %+v, want unavailable with a reason", report.Usage)
		}
	})
}

func TestInspectRuntime_RefusesUntrustedTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(*testing.T, *world)
		wantErr error
	}{
		{name: "runtime never started", arrange: func(*testing.T, *world) {}, wantErr: application.ErrNotFound},
		{
			name:    "execution on another cluster",
			arrange: func(t *testing.T, w *world) { w.up(t); w.ledger.last(t, "demo").ClusterUID = "uid-other" },
			wantErr: application.ErrTargetMismatch,
		},
		{
			name: "namespace replaced",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.kubernetes.edit("Namespace", namespaceName, func(object microruntime.Object) { object.Metadata()["uid"] = "uid-replacement" })
			},
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{
			name: "recorded object replaced",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.kubernetes.edit("Service", "web", func(object microruntime.Object) { object.Metadata()["uid"] = "uid-replacement" })
			},
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{
			name: "foreign quota",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.kubernetes.edit("ResourceQuota", "runtime", func(object microruntime.Object) { delete(object.Metadata(), "labels") })
			},
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{name: "cluster unreachable", arrange: func(t *testing.T, w *world) { w.up(t); w.kubernetes.getErr = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)

			if _, err := w.inspectRuntime().Execute(context.Background(), w.runtimeRequest()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.evidence.inspections) != 0 {
				t.Error("a refused inspection recorded evidence")
			}
		})
	}
}

func TestInspectRuntime_ReportsUnrecordedEvidence(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)
	w.evidence.err = errInjected

	if report, err := w.inspectRuntime().Execute(context.Background(), w.runtimeRequest()); !errors.Is(err, errInjected) || report.ObservationID != "" {
		t.Fatalf("Execute() = %+v, %v; want the evidence failure and no report", report, err)
	}
}
