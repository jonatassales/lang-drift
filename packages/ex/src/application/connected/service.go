package connected

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/exphere/ex/src/domain/connection"
	"github.com/exphere/ex/src/domain/project"
)

func identifier() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// Connect binds one checkout without changing its native project identity. The
// runtime credential is never a command argument or returned to presentation.
func (service Service) Connect(ctx context.Context, request ConnectRequest) (Result, error) {
	if request.Secret == "" {
		return Result{}, errors.New("RUNTIME_SECRET must be supplied in the environment")
	}
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}
	env, err := identifier()
	if err != nil {
		return Result{}, err
	}
	binding, err := service.Console.Bind(ctx, request.URL, request.Secret, env)
	if err != nil {
		return Result{}, err
	}
	if binding.Credential == "" || binding.RuntimeUID == "" || binding.RuntimeID == "" || binding.EnvironmentID != env {
		return Result{}, errors.New("console returned an incomplete binding")
	}
	if err = service.Native.VerifyRepository(ctx, owner, binding.Repository); err != nil {
		return Result{}, err
	}
	existing, err := service.Store.List(owner)
	if err != nil {
		return Result{}, err
	}
	for _, record := range existing {
		if record.RuntimeUID == binding.RuntimeUID {
			return Result{}, errors.New("runtime is already connected in this checkout")
		}
	}
	record := Connection{URL: request.URL, EnvironmentID: env, RuntimeUID: binding.RuntimeUID, RuntimeID: binding.RuntimeID, ProjectID: binding.ProjectID, Repository: binding.Repository, Credential: binding.Credential, Secret: request.Secret, State: "stopped", Managed: map[string]string{}, Outbox: []Event{}}
	if err = service.Store.Bind(ctx, owner, record); err != nil {
		return Result{}, err
	}
	return Result{Runtime: record.RuntimeID, Environment: env, Status: "connected; run ex agent run to synchronize"}, nil
}

// Run owns one worker per connected runtime and joins all workers before exit.
// Each worker maintains communication while finite native operations run.
func (service Service) Run(ctx context.Context, request RunRequest) error {
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
	if len(records) == 0 {
		return errors.New("no connected runtimes; run ex agent connect <console-url>")
	}
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	results := make(chan error, len(records))
	var group sync.WaitGroup
	for _, record := range records {
		group.Add(1)
		go func(record Connection) { defer group.Done(); results <- service.runOne(child, owner, record) }(record)
	}
	go func() { group.Wait(); close(results) }()
	var failures []error
	for err := range results {
		if err != nil {
			failures = append(failures, err)
			cancel()
		}
	}
	return errors.Join(failures...)
}

func (service Service) Migrate(ctx context.Context, request MigrateRequest) (Result, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}
	if request.Recover {
		if err := service.Native.Stop(ctx, owner, Connection{RuntimeID: request.Runtime}); err != nil {
			return Result{}, err
		}
	}
	path, err := service.Native.Migrate(ctx, owner, request.Runtime)
	if err != nil {
		return Result{}, err
	}
	return Result{Runtime: request.Runtime, Path: path, Status: "migration draft written; import the existing runtime ID in the console and supply its build recipes"}, nil
}

func (service Service) interval() time.Duration {
	if service.PollInterval > 0 {
		return service.PollInterval
	}
	return 5 * time.Second
}
func (service Service) timeout() time.Duration {
	if service.OperationTimeout > 0 {
		return service.OperationTimeout
	}
	return 120 * time.Second
}

type workerEvent struct {
	kind, detail string
	ack          chan error
}
type workerResult struct {
	action      string
	cause       string
	operationID string
	err         error
	observation *Observation
}

// runOne owns all connection state. Workers send bounded results back to this
// loop; only this loop persists state, so network polling cannot race the ledger.
func (service Service) runOne(ctx context.Context, owner project.Project, record Connection) error {
	unlock, err := service.Store.Lock(owner, record.RuntimeUID)
	if err != nil {
		return err
	}
	defer unlock()
	var active context.CancelFunc
	var results <-chan workerResult
	reports := make(chan workerEvent)
	lastContact := service.now()
	lease := 15 * time.Second
	var intent Intent
	online := false
	stopCause := ""
	activeCause := ""
	activeOperation := ""
	activeAction := ""
	if record.State != "stopped" {
		record.State = "unknown"
		stopCause = "agent_restart"
	}
	save := func() error { return service.Store.Save(owner, record) }
	appendEvent := func(kind, cause, detail string) error {
		eventID, err := identifier()
		if err != nil {
			return err
		}
		eventRevision := record.Revision
		if eventRevision < 1 || cause == "configuration" {
			eventRevision = intent.Revision
		}
		if eventRevision < 1 {
			eventRevision = 1
		}
		eventExecution := record.ExecutionID
		if cause == "configuration" {
			eventExecution = ""
		}
		event := Event{ID: eventID, Sequence: record.Sequence + 1, PurgeEpoch: record.PurgeEpoch, ExecutionID: eventExecution, Kind: kind, Cause: cause, OccurredAt: time.Now().UTC(), Revision: eventRevision, IntentVersion: record.IntentVersion, Detail: detail, OperationID: record.OperationID, ConcreteHash: record.ConcreteHash, MarkdownHash: markdownHash(record.Markdown)}
		if kind == "build_started" && (len(detail) == 40 || len(detail) == 64) {
			event.SourceRevision = detail
			event.Detail = "build_started"
		}
		if err := addEvent(&record, event); err != nil {
			return err
		}
		return save()
	}
	launch := func(action, cause string, config connection.Configuration) error {
		if active != nil {
			return errors.New("operation already active")
		}
		record.State = map[string]string{"start": "starting", "stop": "stopping", "observe": "running"}[action]
		activeAction = action
		activeCause = cause
		activeOperation = record.OperationID
		if action != "observe" {
			if err := appendEvent(action+"_started", cause, ""); err != nil {
				return err
			}
		}
		workerCtx, cancel := context.WithTimeout(ctx, service.timeout())
		active = cancel
		completed := make(chan workerResult, 1)
		results = completed
		snapshot := record
		go func() {
			result := workerResult{action: action, cause: cause, operationID: snapshot.OperationID}
			report := func(kind, detail string) error {
				ack := make(chan error, 1)
				select {
				case reports <- workerEvent{kind: kind, detail: detail, ack: ack}:
				case <-workerCtx.Done():
					return workerCtx.Err()
				}
				select {
				case err := <-ack:
					return err
				case <-workerCtx.Done():
					return workerCtx.Err()
				}
			}
			switch action {
			case "start":
				result.err = service.Native.Start(workerCtx, owner, snapshot, config, report)
			case "stop":
				result.err = service.Native.Stop(workerCtx, owner, snapshot)
			case "observe":
				result.observation, result.err = service.Native.Observe(workerCtx, owner, snapshot)
			}
			completed <- result
		}()
		return nil
	}
	shutdown := func() error {
		if active != nil {
			active()
			for results != nil {
				select {
				case report := <-reports:
					report.ack <- context.Canceled
				case <-results:
					results = nil
				}
			}
		}
		if record.State == "stopped" {
			return save()
		}
		record.State = "stopping"
		eventErr := appendEvent("stop_started", "agent_shutdown", "")
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), service.timeout())
		defer cancel()
		stopErr := service.Native.Stop(cleanup, owner, record)
		if stopErr != nil {
			record.State = "unknown"
			_ = appendEvent("failed", "agent_shutdown", "operation_failed")
			return errors.Join(eventErr, stopErr)
		}
		record.State = "stopped"
		if err := appendEvent("stopped", "agent_shutdown", ""); err != nil {
			return errors.Join(eventErr, err)
		}
		response, syncErr := service.Console.Sync(cleanup, record, syncOf(record))
		if syncErr == nil {
			record.Outbox = acknowledge(record.Outbox, response.AcceptedEventIDs)
			_ = save()
		}
		return eventErr
	}
	if record.State == "unknown" {
		if err := launch("stop", stopCause, connection.Configuration{}); err != nil {
			return err
		}
	}
	ticker := time.NewTicker(service.interval())
	defer ticker.Stop()
	tick := make(chan time.Time, 1)
	tick <- time.Now()
	for {
		select {
		case <-ctx.Done():
			return shutdown()
		case report := <-reports:
			record.OperationID = activeOperation
			report.ack <- appendEvent(report.kind, activeCause, report.detail)
		case result := <-results:
			active()
			active = nil
			activeAction = ""
			results = nil
			record.OperationID = result.operationID
			switch result.action {
			case "start":
				if result.err != nil {
					record.State = "unknown"
					if stopCause == "" {
						stopCause = "start_failed"
					}
					if err := appendEvent("failed", stopCause, "operation_failed"); err != nil {
						return errors.Join(err, shutdown())
					}
				} else {
					record.State = "running"
					if err := appendEvent("started", result.cause, "native_ready"); err != nil {
						return errors.Join(err, shutdown())
					}
				}
			case "stop":
				if result.err != nil {
					record.State = "unknown"
					if err := appendEvent("failed", result.cause, "operation_failed"); err != nil {
						return errors.Join(err, shutdown())
					}
				} else {
					record.State = "stopped"
					record.Observation = nil
					if err := appendEvent("stopped", result.cause, ""); err != nil {
						return err
					}
					if result.cause == "connection_lost" {
						record.ResumeEligible = true
					}
					stopCause = ""
				}
			case "observe":
				if result.err != nil && stopCause != "" {
					record.Observation = nil
					continue
				}
				if result.err == nil && (result.observation == nil || result.observation.State == "" || result.observation.State == "running") {
					record.Observation = result.observation
				} else {
					record.Observation = nil
					record.State = "unknown"
					stopCause = "runtime_failure"
					if err := appendEvent("failed", stopCause, "operation_failed"); err != nil {
						return errors.Join(err, shutdown())
					}
				}
				if err := save(); err != nil {
					return errors.Join(err, shutdown())
				}
			}
		case <-tick:
			response, syncErr := service.Console.Sync(ctx, record, syncOf(record))
			if syncErr != nil {
				online = false
				if service.now().Sub(lastContact) >= lease && record.State != "stopped" {
					stopCause = "connection_lost"
					if active != nil {
						if activeAction != "stop" {
							active()
						}
					} else if err := launch("stop", stopCause, connection.Configuration{}); err != nil {
						return errors.Join(err, shutdown())
					}
				}
				continue
			}
			online = true
			lastContact = service.now()
			intent = response
			if response.LeaseSeconds >= 5 && response.LeaseSeconds <= 60 {
				lease = time.Duration(response.LeaseSeconds) * time.Second
			}
			record.Outbox = acknowledge(record.Outbox, response.AcceptedEventIDs)
			if response.PurgeEpoch > record.PurgeEpoch {
				record.PurgeEpoch = response.PurgeEpoch
				record.Observation = nil
			}
			if err := save(); err != nil {
				return errors.Join(err, shutdown())
			}
			if pending := response.UnresolvedExecution; pending != nil && pending.ExecutionID != record.ExecutionID {
				if pending.ExecutionID == "" || pending.Revision < 1 || pending.IntentVersion < 1 {
					return errors.Join(errors.New("console returned an invalid reconciliation identity"), shutdown())
				}
				if active != nil {
					if activeAction != "stop" {
						active()
					}
					continue
				}
				// A claim may have committed remotely before its response or local
				// persistence was lost. Adopt identity only to reconcile cleanup.
				record.ExecutionID = pending.ExecutionID
				record.IntentVersion = pending.IntentVersion
				record.Revision = pending.Revision
				record.LastAttempt = pending.IntentVersion
				record.OperationID = response.Operation.ID
				record.State = "unknown"
				record.Observation = nil
				record.ResumeEligible = !response.Revoked && !response.Deleted && response.DesiredState == "running"
				stopCause = "agent_restart"
				if err := save(); err != nil {
					return errors.Join(err, shutdown())
				}
			}
			if response.Revoked || response.Deleted || response.DesiredState != "running" {
				record.ResumeEligible = false
				record.OperationID = response.Operation.ID
				cause := "user_stop"
				if response.Deleted {
					cause = "deleted"
				} else if response.Revoked {
					cause = "revoked"
				}
				if record.State != "stopped" {
					stopCause = cause
					if active != nil {
						if activeAction != "stop" {
							active()
						}
					} else if err := launch("stop", cause, connection.Configuration{}); err != nil {
						return errors.Join(err, shutdown())
					}
				}
			}
			if active != nil {
				continue
			}
			if stopCause != "" && record.State != "stopped" {
				if err := launch("stop", stopCause, connection.Configuration{}); err != nil {
					return errors.Join(err, shutdown())
				}
				continue
			}
			if response.Revoked || response.Deleted || len(record.Outbox) > 0 {
				continue
			}
			if record.State == "stopped" && response.Revision != record.Revision {
				config, err := connection.Parse(response.Markdown)
				if err == nil && (config.RuntimeUID != record.RuntimeUID || config.RuntimeID != record.RuntimeID || config.ProjectID != record.ProjectID || config.Revision != response.Revision) {
					err = errors.New("configuration identities differ from binding")
				}
				if err == nil {
					record.Managed, record.ConcreteHash, err = service.Native.Materialize(ctx, owner, record, config, response.Markdown)
				}
				if err != nil {
					if record.FailureRevision != response.Revision {
						record.FailureRevision = response.Revision
						if err := appendEvent("failed", "configuration", "configuration_conflict"); err != nil {
							return err
						}
					}
					continue
				}
				record.Revision = response.Revision
				record.ExecutionID = ""
				record.OperationID = response.Operation.ID
				record.Markdown = response.Markdown
				if err := appendEvent("configured", "configuration", "configuration_applied"); err != nil {
					return err
				}
				continue // acknowledge materialization before claiming a start
			}
			if response.DesiredState == "running" && record.State == "stopped" && response.Revision == record.Revision && (record.LastAttempt != response.IntentVersion || record.ResumeEligible) {
				config, err := connection.Parse(record.Markdown)
				if err != nil {
					return err
				}
				claim, err := service.Console.Claim(ctx, record, response)
				if err != nil {
					continue
				}
				if claim.Revision != record.Revision || claim.IntentVersion != response.IntentVersion || claim.ExecutionID == "" {
					return errors.New("console returned an inconsistent start claim")
				}
				record.ExecutionID = claim.ExecutionID
				record.IntentVersion = claim.IntentVersion
				record.LastAttempt = claim.IntentVersion
				record.OperationID = response.Operation.ID
				cause := response.Operation.Cause
				if record.ResumeEligible {
					cause = "reconnect"
					record.ResumeEligible = false
				}
				if cause == "" {
					cause = "user_start"
				}
				if err := launch("start", cause, config); err != nil {
					return errors.Join(err, shutdown())
				}
			} else if online && record.State == "running" {
				if err := launch("observe", "observation", connection.Configuration{}); err != nil {
					return errors.Join(err, shutdown())
				}
			}
		case now := <-ticker.C:
			select {
			case tick <- now:
			default:
			}
		}
	}
}

func syncOf(record Connection) Sync {
	events := record.Outbox
	if len(events) > 100 {
		events = events[:100]
	}
	return Sync{EnvironmentID: record.EnvironmentID, Events: events, State: record.State, ExecutionID: record.ExecutionID, LastRevision: record.Revision, PurgeEpoch: record.PurgeEpoch, Observation: record.Observation}
}
func acknowledge(events []Event, ids []string) []Event {
	seen := map[string]bool{}
	for _, id := range ids {
		seen[id] = true
	}
	kept := make([]Event, 0, len(events))
	for _, event := range events {
		if !seen[event.ID] {
			kept = append(kept, event)
		}
	}
	return kept
}

// String does not expose credentials when the connection is logged accidentally.
func (record Connection) String() string {
	return fmt.Sprintf("runtime=%s environment=%s state=%s", record.RuntimeID, record.EnvironmentID, record.State)
}

func markdownHash(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	value := sha256.Sum256([]byte(text))
	return hex.EncodeToString(value[:])
}

func (service Service) now() time.Time {
	if service.Now != nil {
		return service.Now()
	}
	return time.Now()
}
