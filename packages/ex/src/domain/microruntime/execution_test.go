package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

func projectID(t testing.TB) project.ID {
	t.Helper()

	id, err := project.ParseID("fedcba9876543210fedcba9876543210")
	if err != nil {
		t.Fatal(err)
	}

	return id
}

// startedExecution is a new up of validDescriptor's runtime, before preparation.
func startedExecution(t testing.TB) *microruntime.Execution {
	t.Helper()

	return microruntime.StartExecution(microruntime.Start{
		ProjectID:  projectID(t),
		Runtime:    identifier(t, "demo"),
		Namespace:  identifier(t, "demo-local"),
		ClusterUID: "uid-kube-system",
		Context:    "kind-ex-0123456789ab",
		ConfigHash: "hash",
		Token:      token,
		RunID:      "run-up",
		Tools:      map[string]microruntime.ToolVersion{"kubectl": {Version: "v1.32.0", Observed: true}},
		Manifests:  validManifests(t),
	})
}

// ownedNamespace is the live namespace preparation creates for the runtime.
func ownedNamespace(uid string) microruntime.Object {
	return liveObject("v1", "Namespace", "demo-local", uid, map[string]any{microruntime.LabelOwnershipToken: token})
}

// created is the object the API server returns for a claimed resource.
func created(resource microruntime.Resource, uid string) microruntime.Object {
	return liveObject(resource.APIVersion, resource.Kind, resource.Name, uid, map[string]any{microruntime.LabelOwnershipToken: token})
}

// runningExecution is an up that completed: every object was created and
// every unit completed.
func runningExecution(t testing.TB) *microruntime.Execution {
	t.Helper()

	execution := startedExecution(t)
	if err := execution.RecordNamespace(ownedNamespace("uid-namespace")); err != nil {
		t.Fatal(err)
	}

	if err := execution.BeginApplying(); err != nil {
		t.Fatal(err)
	}

	for index, resource := range execution.Resources {
		if err := execution.RecordCreated(index, created(resource, "uid-"+resource.Kind+"-"+resource.Name)); err != nil {
			t.Fatal(err)
		}

		execution.MarkUnit(resource.Unit, microruntime.UnitCompleted)
	}

	if err := execution.CompleteUp(); err != nil {
		t.Fatal(err)
	}

	execution.Settle(nil, false)

	return execution
}

func TestStartExecution(t *testing.T) {
	t.Parallel()

	manifests := validManifests(t)
	execution := microruntime.StartExecution(microruntime.Start{
		ProjectID: projectID(t), Runtime: identifier(t, "demo"), Namespace: identifier(t, "demo-local"),
		ClusterUID: "uid-kube-system", Context: "kind-ex", ConfigHash: "hash", Token: token, RunID: "run-up",
		Manifests: manifests,
	})

	if execution.Status != microruntime.StatusPreparing || execution.Operation != microruntime.OperationUp ||
		execution.Settled || execution.NamespaceUID != "" {
		t.Errorf("StartExecution() = %+v, want an unsettled up that is preparing", execution)
	}

	if len(execution.Resources) != len(manifests) {
		t.Fatalf("recorded %d resources, want one per manifest", len(execution.Resources))
	}

	for index, resource := range execution.Resources {
		object := manifests[index].Object
		if resource.Kind != object.Kind() || resource.Name != object.Name() || resource.Namespace != "demo-local" ||
			resource.UID != "" || resource.Unit != manifests[index].Unit {
			t.Errorf("resource %d = %+v, want the claimed manifest without a UID", index, resource)
		}

		if object.Labels().Text(microruntime.LabelOwnershipToken) != token {
			t.Errorf("manifest %s was not claimed", manifests[index].Key())
		}
	}

	if len(execution.Units) != 2 {
		t.Fatalf("units = %v, want web and batch", execution.Units)
	}

	for unit, status := range execution.Units {
		if status != microruntime.UnitNotAttempted {
			t.Errorf("unit %s = %s, want not-attempted", unit, status)
		}
	}

	if len(execution.Images) != 2 {
		t.Errorf("images = %v, want the two configured containers", execution.Images)
	}
}

func TestExecution_Validate(t *testing.T) {
	t.Parallel()

	runtime := identifier(t, "demo")
	other := identifier(t, "other")

	tests := []struct {
		name    string
		edit    func(*microruntime.Execution)
		runtime microruntime.Identifier
		wantErr bool
	}{
		{name: "consistent record", edit: func(*microruntime.Execution) {}},
		{name: "record of another runtime", edit: func(*microruntime.Execution) {}, runtime: other, wantErr: true},
		{name: "record of another project", edit: func(e *microruntime.Execution) { e.ProjectID = project.ID{} }, wantErr: true},
		{name: "missing token", edit: func(e *microruntime.Execution) { e.Token = "" }, wantErr: true},
		{name: "missing namespace", edit: func(e *microruntime.Execution) { e.Namespace = microruntime.Identifier{} }, wantErr: true},
		{name: "missing cluster UID", edit: func(e *microruntime.Execution) { e.ClusterUID = "" }, wantErr: true},
		{name: "unsupported kind", edit: func(e *microruntime.Execution) { e.Resources[0].Kind = "Secret" }, wantErr: true},
		{name: "mismatched API version", edit: func(e *microruntime.Execution) { e.Resources[1].APIVersion = "apps/v1" }, wantErr: true},
		{name: "invalid name", edit: func(e *microruntime.Execution) { e.Resources[0].Name = "../web" }, wantErr: true},
		{name: "another namespace", edit: func(e *microruntime.Execution) { e.Resources[0].Namespace = "kube-system" }, wantErr: true},
		{name: "resource without a unit", edit: func(e *microruntime.Execution) { e.Resources[2].Unit = microruntime.Identifier{} }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := runningExecution(t)
			tt.edit(execution)

			candidate := runtime
			if !tt.runtime.IsZero() {
				candidate = tt.runtime
			}

			err := execution.Validate(projectID(t), candidate)

			if tt.wantErr != errors.Is(err, microruntime.ErrInconsistentExecution) || (!tt.wantErr && err != nil) {
				t.Fatalf("Validate() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestExecution_BeginApplying(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		edit    func(*microruntime.Execution)
		wantErr bool
	}{
		{name: "prepared startup", edit: func(e *microruntime.Execution) { e.NamespaceUID = "uid-namespace" }},
		{name: "no recorded namespace", edit: func(*microruntime.Execution) {}, wantErr: true},
		{
			name:    "settled startup",
			edit:    func(e *microruntime.Execution) { e.NamespaceUID, e.Settled = "uid-namespace", true },
			wantErr: true,
		},
		{
			name: "cleanup",
			edit: func(e *microruntime.Execution) {
				e.NamespaceUID = "uid-namespace"
				e.Operation = microruntime.OperationDown
			},
			wantErr: true,
		},
		{
			name: "already applying",
			edit: func(e *microruntime.Execution) {
				e.NamespaceUID = "uid-namespace"
				e.Status = microruntime.StatusApplying
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := startedExecution(t)
			tt.edit(execution)
			before := execution.Status

			err := execution.BeginApplying()

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrInconsistentExecution) || execution.Status != before {
					t.Fatalf("BeginApplying() = %v with status %s, want refusal without a transition", err, execution.Status)
				}

				return
			}

			if err != nil || execution.Status != microruntime.StatusApplying {
				t.Fatalf("BeginApplying() = %v with status %s, want applying", err, execution.Status)
			}
		})
	}
}

func TestExecution_CompleteUp(t *testing.T) {
	t.Parallel()

	// applied returns an up whose objects were all created and whose units
	// all completed, so each case breaks exactly one condition.
	applied := func(t *testing.T) *microruntime.Execution {
		execution := startedExecution(t)
		execution.NamespaceUID = "uid-namespace"
		execution.Status = microruntime.StatusApplying

		for index := range execution.Resources {
			execution.Resources[index].UID = "uid-" + execution.Resources[index].Name
			execution.MarkUnit(execution.Resources[index].Unit, microruntime.UnitCompleted)
		}

		return execution
	}
	extra := identifier(t, "extra")

	tests := []struct {
		name    string
		edit    func(*microruntime.Execution)
		wantErr bool
	}{
		{name: "every resource created and unit completed", edit: func(*microruntime.Execution) {}},
		{name: "resource without a UID", edit: func(e *microruntime.Execution) { e.Resources[1].UID = "" }, wantErr: true},
		{
			name:    "unit still running",
			edit:    func(e *microruntime.Execution) { e.MarkUnit(e.Resources[2].Unit, microruntime.UnitRunning) },
			wantErr: true,
		},
		{
			name:    "recorded unit without resources",
			edit:    func(e *microruntime.Execution) { e.Units[extra] = microruntime.UnitNotAttempted },
			wantErr: true,
		},
		{name: "no resources", edit: func(e *microruntime.Execution) { e.Resources = nil }, wantErr: true},
		{name: "no namespace", edit: func(e *microruntime.Execution) { e.NamespaceUID = "" }, wantErr: true},
		{name: "still preparing", edit: func(e *microruntime.Execution) { e.Status = microruntime.StatusPreparing }, wantErr: true},
		{name: "settled", edit: func(e *microruntime.Execution) { e.Settled = true }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := applied(t)
			tt.edit(execution)

			err := execution.CompleteUp()

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrInconsistentExecution) || execution.Status == microruntime.StatusRunning {
					t.Fatalf("CompleteUp() = %v with status %s, want refusal", err, execution.Status)
				}

				return
			}

			if err != nil || execution.Status != microruntime.StatusRunning {
				t.Fatalf("CompleteUp() = %v with status %s, want running", err, execution.Status)
			}
		})
	}
}

func TestExecution_RequireClosed(t *testing.T) {
	t.Parallel()

	closed := func(t *testing.T) *microruntime.Execution {
		execution := runningExecution(t)
		execution.BeginDown("run-down", nil)
		execution.CompleteDown()
		execution.Settle(nil, false)

		return execution
	}

	tests := []struct {
		name    string
		edit    func(*microruntime.Execution)
		wantErr bool
	}{
		{name: "settled cleanup", edit: func(*microruntime.Execution) {}},
		{name: "unsettled operation", edit: func(e *microruntime.Execution) { e.Settled = false }, wantErr: true},
		{name: "running runtime", edit: func(e *microruntime.Execution) { e.Status = microruntime.StatusRunning }, wantErr: true},
		{name: "failed cleanup", edit: func(e *microruntime.Execution) { e.Status = microruntime.StatusFailed }, wantErr: true},
		{name: "remaining namespace", edit: func(e *microruntime.Execution) { e.NamespaceUID = "uid-namespace" }, wantErr: true},
		{
			name:    "remaining resource",
			edit:    func(e *microruntime.Execution) { e.Resources = []microruntime.Resource{{Kind: "Service"}} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := closed(t)
			tt.edit(execution)

			err := execution.RequireClosed()

			if tt.wantErr != errors.Is(err, microruntime.ErrExistingObligations) || (!tt.wantErr && err != nil) {
				t.Fatalf("RequireClosed() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestExecution_AuthorizeDown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  microruntime.ExecutionStatus
		settled bool
		recover bool
		wantErr bool
	}{
		{name: "running runtime", status: microruntime.StatusRunning, settled: true},
		{name: "running runtime with recovery", status: microruntime.StatusRunning, settled: true, recover: true},
		{name: "failed startup without recovery", status: microruntime.StatusFailed, settled: true, wantErr: true},
		{name: "failed startup with recovery", status: microruntime.StatusFailed, settled: true, recover: true},
		{name: "interrupted startup with recovery", status: microruntime.StatusInterrupted, settled: true, recover: true},
		{name: "unsettled operation", status: microruntime.StatusRunning, wantErr: true},
		{name: "unsettled operation with recovery", status: microruntime.StatusApplying, recover: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := &microruntime.Execution{Runtime: identifier(t, "demo"), Status: tt.status, Settled: tt.settled}
			err := execution.AuthorizeDown(tt.recover)

			if tt.wantErr != errors.Is(err, microruntime.ErrCleanupBlocked) || (!tt.wantErr && err != nil) {
				t.Fatalf("AuthorizeDown(%t) error = %v, want refusal %t", tt.recover, err, tt.wantErr)
			}
		})
	}
}

func TestExecution_DownLifecycle(t *testing.T) {
	t.Parallel()

	execution := runningExecution(t)
	tools := map[string]microruntime.ToolVersion{"kind": {Observed: false}}

	execution.BeginDown("run-down", tools)

	if execution.RunID != "run-down" || execution.Operation != microruntime.OperationDown ||
		execution.Status != microruntime.StatusCleaning || execution.Settled || execution.ToolVersions["kind"].Observed {
		t.Fatalf("BeginDown() = %+v, want an unsettled cleanup run", execution)
	}

	for unit, status := range execution.Units {
		if status != microruntime.UnitNotAttempted {
			t.Errorf("unit %s = %s after BeginDown, want not-attempted", unit, status)
		}
	}

	execution.CompleteDown()
	execution.Settle(nil, false)

	if execution.Status != microruntime.StatusDown || !execution.Settled ||
		execution.NamespaceUID != "" || execution.Resources != nil {
		t.Fatalf("CompleteDown() = %+v, want a settled record without obligations", execution)
	}

	for unit, status := range execution.Units {
		if status != microruntime.UnitCompleted {
			t.Errorf("unit %s = %s after CompleteDown, want completed", unit, status)
		}
	}

	if err := execution.RequireClosed(); err != nil {
		t.Fatalf("RequireClosed() after a completed cleanup = %v", err)
	}
}

func TestExecution_Settle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		failure     error
		interrupted bool
		wantStatus  microruntime.ExecutionStatus
		wantRunning microruntime.UnitStatus
	}{
		{name: "success keeps the reached status", wantStatus: microruntime.StatusApplying, wantRunning: microruntime.UnitRunning},
		{name: "failure", failure: errors.New("dry run failed"), wantStatus: microruntime.StatusFailed, wantRunning: microruntime.UnitFailed},
		{
			name: "interruption", failure: errors.New("canceled"), interrupted: true,
			wantStatus: microruntime.StatusInterrupted, wantRunning: microruntime.UnitFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := startedExecution(t)
			execution.Status = microruntime.StatusApplying
			web, batch := identifier(t, "web"), identifier(t, "batch")
			execution.MarkUnit(web, microruntime.UnitRunning)
			execution.MarkUnit(batch, microruntime.UnitCompleted)

			execution.Settle(tt.failure, tt.interrupted)

			if !execution.Settled || execution.Status != tt.wantStatus {
				t.Fatalf("Settle() = status %s, settled %t; want %s, true", execution.Status, execution.Settled, tt.wantStatus)
			}

			if execution.Units[web] != tt.wantRunning || execution.Units[batch] != microruntime.UnitCompleted {
				t.Fatalf("units = %v, want running units %s and completed units kept", execution.Units, tt.wantRunning)
			}
		})
	}
}

func TestExecution_RecordNamespace(t *testing.T) {
	t.Parallel()

	t.Run("records the owned namespace", func(t *testing.T) {
		t.Parallel()

		execution := startedExecution(t)
		if err := execution.RecordNamespace(ownedNamespace("uid-namespace")); err != nil || execution.NamespaceUID != "uid-namespace" {
			t.Fatalf("RecordNamespace() = %v, NamespaceUID %q", err, execution.NamespaceUID)
		}
	})

	for name, live := range map[string]microruntime.Object{
		"absent namespace":     nil,
		"unlabelled namespace": liveObject("v1", "Namespace", "demo-local", "uid-namespace", nil),
		"another token":        liveObject("v1", "Namespace", "demo-local", "uid-namespace", map[string]any{microruntime.LabelOwnershipToken: "other"}),
	} {
		t.Run("refuses: "+name, func(t *testing.T) {
			t.Parallel()

			if err := startedExecution(t).RecordNamespace(live); !errors.Is(err, microruntime.ErrOwnershipChanged) {
				t.Fatalf("RecordNamespace() error = %v, want ErrOwnershipChanged", err)
			}
		})
	}
}

func TestExecution_ConfirmNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		recorded string
		live     microruntime.Object
		wantUID  string
		wantErr  bool
	}{
		{name: "recorded namespace", recorded: "uid-namespace", live: ownedNamespace("uid-namespace"), wantUID: "uid-namespace"},
		{name: "learns the UID a partial up never recorded", live: ownedNamespace("uid-namespace"), wantUID: "uid-namespace"},
		{name: "replaced namespace", recorded: "uid-namespace", live: ownedNamespace("uid-replacement"), wantErr: true},
		{
			name:    "foreign namespace a partial up never recorded",
			live:    liveObject("v1", "Namespace", "demo-local", "uid-foreign", nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := startedExecution(t)
			execution.NamespaceUID = tt.recorded

			err := execution.ConfirmNamespace(tt.live)

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrOwnershipChanged) {
					t.Fatalf("ConfirmNamespace() error = %v, want ErrOwnershipChanged", err)
				}

				return
			}

			if err != nil || execution.NamespaceUID != tt.wantUID {
				t.Fatalf("ConfirmNamespace() = %v, NamespaceUID %q; want %q", err, execution.NamespaceUID, tt.wantUID)
			}
		})
	}
}

func TestExecution_RecordCreated(t *testing.T) {
	t.Parallel()

	t.Run("records the identity the API server returned", func(t *testing.T) {
		t.Parallel()

		execution := startedExecution(t)
		resource := execution.Resources[0]

		if err := execution.RecordCreated(0, created(resource, "uid-created")); err != nil {
			t.Fatal(err)
		}

		if got := execution.Resources[0]; got.UID != "uid-created" || got.Unit != resource.Unit || got.Name != resource.Name {
			t.Fatalf("recorded %+v, want the created UID with the original unit", got)
		}
	})

	t.Run("refuses an object without the runtime's ownership", func(t *testing.T) {
		t.Parallel()

		execution := startedExecution(t)
		resource := execution.Resources[0]
		unclaimed := liveObject(resource.APIVersion, resource.Kind, resource.Name, "uid-created", nil)

		if err := execution.RecordCreated(0, unclaimed); !errors.Is(err, microruntime.ErrOwnershipChanged) {
			t.Fatalf("RecordCreated() error = %v, want ErrOwnershipChanged", err)
		}
	})
}

func TestExecution_ConfirmResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		recorded string
		live     func(microruntime.Resource) microruntime.Object
		wantUID  string
		wantErr  bool
	}{
		{name: "already removed", recorded: "uid-web", live: func(microruntime.Resource) microruntime.Object { return nil }, wantUID: "uid-web"},
		{name: "recorded object", recorded: "uid-web", live: func(r microruntime.Resource) microruntime.Object { return created(r, "uid-web") }, wantUID: "uid-web"},
		{
			name:    "learns the UID a partial up never recorded",
			live:    func(r microruntime.Resource) microruntime.Object { return created(r, "uid-web") },
			wantUID: "uid-web",
		},
		{
			name:     "replacement under the same name",
			recorded: "uid-web",
			live:     func(r microruntime.Resource) microruntime.Object { return created(r, "uid-replacement") },
			wantErr:  true,
		},
		{
			name: "foreign object a partial up never recorded",
			live: func(r microruntime.Resource) microruntime.Object {
				return liveObject(r.APIVersion, r.Kind, r.Name, "uid-foreign", nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			execution := startedExecution(t)
			execution.Resources[0].UID = tt.recorded

			err := execution.ConfirmResource(0, tt.live(execution.Resources[0]))

			if tt.wantErr {
				if !errors.Is(err, microruntime.ErrOwnershipChanged) {
					t.Fatalf("ConfirmResource() error = %v, want ErrOwnershipChanged", err)
				}

				return
			}

			if err != nil || execution.Resources[0].UID != tt.wantUID {
				t.Fatalf("ConfirmResource() = %v, UID %q; want %q", err, execution.Resources[0].UID, tt.wantUID)
			}
		})
	}
}

func TestExecution_Progress(t *testing.T) {
	t.Parallel()

	execution := startedExecution(t)
	progress := execution.Progress()
	web, batch := identifier(t, "web"), identifier(t, "batch")

	if progress.Finish(web) {
		t.Fatal("web finished after one of its two resources")
	}

	if !progress.Finish(web) {
		t.Fatal("web did not finish after its last resource")
	}

	if !progress.Finish(batch) {
		t.Fatal("batch did not finish after its only resource")
	}
}
