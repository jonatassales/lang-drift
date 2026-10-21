package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func TestResourceOf(t *testing.T) {
	t.Parallel()

	unit := identifier(t, "web")
	live := liveObject("apps/v1", "Deployment", "web", "uid-web", nil)

	got := microruntime.ResourceOf(live, unit)
	want := microruntime.Resource{APIVersion: "apps/v1", Kind: "Deployment", Name: "web", Namespace: "demo-local", UID: "uid-web", Unit: unit}

	if got != want {
		t.Fatalf("ResourceOf() = %+v, want %+v", got, want)
	}
}

func TestResource_IsWorkload(t *testing.T) {
	t.Parallel()

	for kind, want := range map[string]bool{"Deployment": true, "Job": true, "Service": false, "Pod": false} {
		if got := (microruntime.Resource{Kind: kind}).IsWorkload(); got != want {
			t.Errorf("IsWorkload() for %s = %t, want %t", kind, got, want)
		}
	}
}

func TestResource_Readiness(t *testing.T) {
	t.Parallel()

	withStatus := func(kind string, spec, status map[string]any, generation float64) microruntime.Object {
		object := liveObject("apps/v1", kind, "web", "uid-web", nil)
		object["spec"], object["status"] = spec, status
		object["metadata"].(map[string]any)["generation"] = generation

		return object
	}
	condition := func(kind, status string) map[string]any {
		return map[string]any{"conditions": []any{map[string]any{"type": kind, "status": status, "reason": "BackoffLimitExceeded"}}}
	}
	settled := map[string]any{"observedGeneration": 2.0, "updatedReplicas": 2.0, "availableReplicas": 2.0, "replicas": 2.0}

	tests := []struct {
		name      string
		kind      string
		live      microruntime.Object
		wantReady bool
		wantErr   error
	}{
		{name: "absent object", kind: "Deployment", live: nil, wantErr: microruntime.ErrWorkloadFailed},
		{
			name:    "replacement under the same name",
			kind:    "Service",
			live:    liveObject("v1", "Service", "web", "uid-other", nil),
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{name: "existing Service", kind: "Service", live: liveObject("v1", "Service", "web", "uid-web", nil), wantReady: true},
		{name: "completed Job", kind: "Job", live: withStatus("Job", nil, condition("Complete", "True"), 1), wantReady: true},
		{name: "failed Job", kind: "Job", live: withStatus("Job", nil, condition("Failed", "True"), 1), wantErr: microruntime.ErrWorkloadFailed},
		{name: "Job condition not yet true", kind: "Job", live: withStatus("Job", nil, condition("Complete", "False"), 1)},
		{name: "running Job", kind: "Job", live: withStatus("Job", nil, map[string]any{"active": 1.0}, 1)},
		{
			name:      "Deployment with every replica available",
			kind:      "Deployment",
			live:      withStatus("Deployment", map[string]any{"replicas": 2.0}, settled, 2),
			wantReady: true,
		},
		{
			name:      "Deployment defaulting to one replica",
			kind:      "Deployment",
			live:      withStatus("Deployment", map[string]any{}, map[string]any{"observedGeneration": 1.0, "updatedReplicas": 1.0, "availableReplicas": 1.0, "replicas": 1.0}, 1),
			wantReady: true,
		},
		{
			name: "Deployment whose status predates its spec",
			kind: "Deployment",
			live: withStatus("Deployment", map[string]any{"replicas": 2.0}, settled, 3),
		},
		{
			name: "Deployment with unavailable replicas",
			kind: "Deployment",
			live: withStatus("Deployment", map[string]any{"replicas": 2.0},
				map[string]any{"observedGeneration": 2.0, "updatedReplicas": 2.0, "availableReplicas": 1.0, "replicas": 2.0}, 2),
		},
		{
			name: "Deployment still terminating old replicas",
			kind: "Deployment",
			live: withStatus("Deployment", map[string]any{"replicas": 2.0},
				map[string]any{"observedGeneration": 2.0, "updatedReplicas": 2.0, "availableReplicas": 2.0, "replicas": 3.0}, 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource := microruntime.Resource{Kind: tt.kind, Name: "web", UID: "uid-web"}
			ready, err := resource.Readiness(tt.live)

			if !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && err != nil) {
				t.Fatalf("Readiness() error = %v, want %v", err, tt.wantErr)
			}

			if ready != tt.wantReady {
				t.Fatalf("Readiness() = %t, want %t", ready, tt.wantReady)
			}
		})
	}
}

func TestRemoved(t *testing.T) {
	t.Parallel()

	if removed, err := microruntime.Removed(nil, "uid-web"); !removed || err != nil {
		t.Errorf("Removed(nil) = %t, %v; want true", removed, err)
	}

	if removed, err := microruntime.Removed(liveObject("v1", "Service", "web", "uid-web", nil), "uid-web"); removed || err != nil {
		t.Errorf("Removed(still present) = %t, %v; want false", removed, err)
	}

	_, err := microruntime.Removed(liveObject("v1", "Service", "web", "uid-other", nil), "uid-web")
	if !errors.Is(err, microruntime.ErrOwnershipChanged) {
		t.Errorf("Removed(replacement) error = %v, want ErrOwnershipChanged", err)
	}
}
