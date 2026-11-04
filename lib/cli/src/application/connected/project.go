package connected

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// RunProject owns one project lock and one worker per microruntime. Discovery
// never grants a start: each worker still requires current intent and a claim.
// Cancellation joins every worker, allowing cleanup and durable event flushing.
func (service Service) RunProject(ctx context.Context, request ProjectRunRequest) error {
	if request.Secret == "" || service.ProjectConsole == nil {
		return errors.New("project agent requires PROJECT_SECRET and a project console")
	}
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return err
	}
	unlock, err := service.Store.Lock(owner, "project-agent")
	if err != nil {
		return err
	}
	defer unlock()
	records, err := service.Store.List(owner)
	if err != nil {
		return err
	}
	child, cancel := context.WithCancel(ctx)
	var workers sync.WaitGroup
	var revoked atomic.Bool
	workerService := service
	workerService.Console = projectRecoveryConsole{Console: service.Console, revoked: &revoked}
	defer func() { cancel(); workers.Wait() }()
	// Bound worker admission and result buffering together; no detached worker
	// can block reporting while the owner joins it during shutdown.
	const capacity = 128
	results := make(chan error, capacity)
	started := make(map[string]bool)
	start := func(record Connection) error {
		if started[record.RuntimeUID] {
			return nil
		}
		if len(started) >= capacity {
			return errors.New("project exceeds 128 connected microruntimes")
		}
		started[record.RuntimeUID] = true
		workers.Add(1)
		go func() { defer workers.Done(); results <- workerService.runOne(child, owner, record) }()
		return nil
	}
	var projectID string
	membership, syncErr := service.ProjectConsole.ProjectSync(child, request.URL, request.Secret)
	revoked.Store(errors.Is(syncErr, ErrAccessRevoked))
	if syncErr == nil {
		projectID = membership.ProjectID
	}
	// Recovery workers must run even when the project credential or network is
	// unavailable after a restart. Their environment credentials only reconcile
	// current console intent and cleanup; revoked access still denies claims.
	for _, record := range records {
		if record.URL != request.URL || (projectID != "" && record.ProjectID != projectID) {
			return errors.New("checkout contains a different console project binding")
		}
		if projectID == "" {
			projectID = record.ProjectID
		}
		if syncErr == nil && membership.Repository.Kind == "local" && record.Repository.Kind != "local" {
			if record.State != "stopped" {
				return errors.New("legacy connection must be stopped before local source association")
			}
			record.Repository = membership.Repository
			if err := service.Store.Save(owner, record); err != nil {
				return err
			}
		}
		if err := start(record); err != nil {
			return err
		}
	}
	ticker := time.NewTicker(service.interval())
	defer ticker.Stop()
	for {
		// Keep existing environment workers alive to reconcile cleanup and flush
		// outboxes. Revoked project access admits no new bindings or starts.
		if errors.Is(syncErr, ErrAccessRevoked) && len(started) == 0 {
			return syncErr
		}
		if syncErr == nil {
			if membership.ProjectID == "" || (projectID != "" && membership.ProjectID != projectID) {
				return errors.New("project membership identity changed")
			}
			projectID = membership.ProjectID
			for _, runtimeID := range membership.RuntimeIDs {
				if started[runtimeID] {
					continue
				}
				if len(started) >= capacity {
					return errors.New("project exceeds 128 connected microruntimes")
				}
				environmentID, err := identifier()
				if err != nil {
					return err
				}
				binding, err := service.ProjectConsole.ProjectBind(child, request.URL, request.Secret, runtimeID, environmentID)
				if err != nil {
					return err
				}
				if binding.ProjectID != projectID || binding.RuntimeUID != runtimeID || binding.RuntimeID == "" || binding.Credential == "" || binding.EnvironmentID != environmentID {
					return errors.New("console returned an inconsistent project binding")
				}
				if err := service.Native.VerifyRepository(child, owner, binding.Repository); err != nil {
					return err
				}
				record := Connection{URL: request.URL, EnvironmentID: environmentID, RuntimeUID: runtimeID, RuntimeID: binding.RuntimeID, ProjectID: projectID, Repository: binding.Repository, Credential: binding.Credential, State: "stopped", Managed: map[string]string{}, Outbox: []Event{}}
				if err := service.Store.Bind(child, owner, record); err != nil {
					return err
				}
				if err := start(record); err != nil {
					return err
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case err := <-results:
			// A cancelled worker and the owner's cancellation can become ready
			// together. Normal shutdown must not become an unexpected-exit error.
			if ctx.Err() != nil && (err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				return nil
			}
			if err == nil {
				return errors.New("microruntime worker exited unexpectedly")
			}
			return err
		case <-ticker.C:
		}
		membership, syncErr = service.ProjectConsole.ProjectSync(child, request.URL, request.Secret)
		if errors.Is(syncErr, ErrAccessRevoked) {
			revoked.Store(true)
		}
	}
}

// Project authority also gates migrated workers whose legacy environment
// credentials predate project secrets. Recovery may sync/stop, never claim.
type projectRecoveryConsole struct {
	Console
	revoked *atomic.Bool
}

func (console projectRecoveryConsole) Sync(ctx context.Context, record Connection, sync Sync) (Intent, error) {
	intent, err := console.Console.Sync(ctx, record, sync)
	if console.revoked.Load() {
		intent.Revoked = true
	}
	return intent, err
}
func (console projectRecoveryConsole) Claim(ctx context.Context, record Connection, intent Intent) (Claim, error) {
	if console.revoked.Load() {
		return Claim{}, ErrAccessRevoked
	}
	return console.Console.Claim(ctx, record, intent)
}
