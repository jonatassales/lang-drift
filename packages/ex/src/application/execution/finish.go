package execution

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// The outcomes a run may end with.
const (
	ReadyForReview = "readyForReview"
	ReadyToMerge   = "readyToMerge"
)

// The check results as the run records them.
const (
	CheckPassed  = "passed"
	CheckFailed  = "failed"
	CheckPending = "pending"
)

// finish is the run after its stages: the runtime comes up, the
// implementation is complete, delivery is observed, the runtime goes down,
// the reservation settles and the run remembers. Each part plays on from
// where an earlier run of the lap left it.
func (playing *playing) finish(ctx context.Context) (StartResult, error) {
	if err := playing.bringUp(ctx); err != nil {
		return playing.result(), err
	}

	playing.completed()

	if err := playing.observeDelivery(ctx); err != nil {
		return playing.result(), err
	}

	if err := playing.takeDown(ctx); err != nil {
		return playing.result(), err
	}

	return playing.release()
}

// bringUp materializes the runtime from its fence, builds and starts it,
// and prints that it is ready, with what the plan proposed for it.
func (playing *playing) bringUp(ctx context.Context) error {
	if playing.lap.Ready != nil {
		return nil
	}

	if playing.lap.Execution == "" {
		if err := playing.start(ctx); err != nil {
			return err
		}
	}

	// A lap that started but was not ready looks again before it goes on.
	observation, err := playing.service.Lifecycle.Observe(ctx, playing.owner, playing.runtime)
	if err != nil {
		return playing.abort(fmt.Errorf("runtime %s could not be observed: %w", playing.runtime, err))
	}
	if !observation.Ready {
		return playing.abort(fmt.Errorf("runtime %s started but is not ready: %d of %d containers ready", playing.runtime, observation.ReadyContainers, observation.Containers))
	}

	now := playing.now()
	playing.lap.Ready = &now
	playing.lap.Resources = &ResourceObservation{Proposed: playing.planned.Totals.Proposed, Executed: harness.Resources{CPUMillis: playing.config.CPUMillis, MemoryMiB: playing.config.MemoryMiB}}
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Ready", At: now, Note: fmt.Sprintf("%d of %d containers ready on %s", observation.ReadyContainers, observation.Containers, playing.readinessSignals())})
	playing.sample(observation)

	emit(playing.sink, session.RuntimeReady, session.Data{"runtime": playing.runtime.String(), "revision": playing.config.Revision, "deployments": observation.Deployments, "cpuMillis": playing.planned.Totals.Proposed.CPUMillis, "memoryMiB": playing.planned.Totals.Proposed.MemoryMiB})

	return playing.persist()
}

// start materializes, builds and starts the runtime, and claims its execution.
func (playing *playing) start(ctx context.Context) error {
	markdown, err := playing.service.Records.ReadRuntime(playing.owner, playing.runtime)
	if err != nil {
		return err
	}

	report := func(step, note string) {
		playing.lap.Steps = append(playing.lap.Steps, Step{Name: step, At: playing.now(), Note: note})
	}

	execution, err := playing.service.Lifecycle.Up(ctx, playing.owner, playing.runtime, playing.config, markdown, report)
	if err != nil {
		return playing.abort(fmt.Errorf("runtime %s did not come up: %w", playing.runtime, err))
	}

	playing.lap.Execution = execution.ID
	playing.lap.Runtime = &run.RuntimeFence{ID: playing.runtime.String(), Revision: playing.config.Revision, MarkdownSHA256: collection.SHA256(markdown)}
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Claimed", At: playing.now(), Note: "execution `" + execution.ID + "`"})

	return nil
}

// readinessSignals names what each component answers readiness with.
func (playing *playing) readinessSignals() string {
	signals := make([]string, 0, len(playing.config.Components))
	for _, component := range playing.config.Components {
		signals = append(signals, "`"+readinessSignal(component)+"`")
	}

	return strings.Join(signals, ", ")
}

// sample keeps the resource peaks and the first steady reading of the
// runtime, and how they were measured.
func (playing *playing) sample(observation Observation) {
	resources := playing.lap.Resources
	if resources == nil {
		return
	}

	resources.Deployments = max(resources.Deployments, observation.Deployments)
	if observation.CPUMillis != nil {
		cpu := *observation.CPUMillis
		if resources.CPUSteady == nil {
			resources.CPUSteady = &cpu
		}
		resources.CPUPeak = number(max(valueOf(resources.CPUPeak), cpu))
	}
	if observation.MemoryMiB != nil {
		resources.MemoryPeak = number(max(valueOf(resources.MemoryPeak), *observation.MemoryMiB))
	}
	if observation.Window != "" {
		resources.Window = observation.Window
	}
	if !observation.ObservedAt.IsZero() && (resources.OldestSample == nil || observation.ObservedAt.Before(resources.OldestSample.Time)) {
		oldest := harness.At(observation.ObservedAt)
		resources.OldestSample = &oldest
	}
}

// completed prints that the implementation is complete: every stage the
// run plays is done, with what the lap observed against its reservation.
func (playing *playing) completed() {
	if playing.lap.StagesCompleted == nil {
		now := playing.now()
		playing.lap.StagesCompleted = &now
		playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Completed", At: now, Note: finops.Dollars(playing.lap.ObservedCents) + " observed of " + finops.Dollars(playing.lap.Reserved) + " reserved"})
	}

	emit(playing.sink, session.Completed, session.Data{"observedCents": playing.lap.ObservedCents, "reservedCents": playing.lap.Reserved})
}

// observeDelivery reads what the harness delivered: the branch and whether
// it was pushed, the pull request and its checks. Nothing it cannot see
// counts as anything; it prints why.
func (playing *playing) observeDelivery(ctx context.Context) error {
	if playing.lap.Delivered != nil {
		return nil
	}

	root := playing.owner.Root
	branch := playing.service.Delivery.Branch(ctx, root)
	if branch.Name != "" {
		playing.record.Delivery.Branch = text(branch.Name)
	}
	if branch.Commit != "" {
		playing.record.Delivery.Commit = text(branch.Commit)
	}

	playing.pushed = branch.Upstream
	if branch.Upstream {
		emit(playing.sink, session.Pushed, session.Data{"branch": branch.Name, "commit": branch.Commit})
	} else {
		emit(playing.sink, session.NotPushed, session.Data{"reason": branch.Reason})
	}

	pullRequest, reason := (*run.PullRequest)(nil), "no branch"
	if branch.Name != "" {
		pullRequest, reason = playing.service.Delivery.PullRequest(ctx, root, branch.Name)
	}

	if pullRequest == nil {
		emit(playing.sink, session.PullRequestUnknown, session.Data{"reason": reason})
		emit(playing.sink, session.ChecksUnknown, session.Data{"reason": "no pull request"})
	} else {
		playing.record.Delivery.PullRequest = pullRequest
		emit(playing.sink, session.PullRequest, session.Data{"number": pullRequest.Number, "title": pullRequest.Title})
		playing.checks(ctx, pullRequest.Number)
	}

	now := playing.now()
	playing.lap.Delivered = &now
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Delivered", At: now, Note: playing.deliveryWords(branch)})

	return playing.persist()
}

// checks follows the pull request's checks until they end or the wait is
// over; a check still running when the wait ends stays pending.
func (playing *playing) checks(ctx context.Context, number int) {
	// The wait is real time, not the records' clock: the checks run elsewhere.
	deadline := time.Now().Add(playing.checksWait)
	announced := false

	for {
		checks, reason := playing.service.Delivery.Checks(ctx, playing.owner.Root, number)
		if reason != "" {
			emit(playing.sink, session.ChecksUnknown, session.Data{"reason": reason})

			return
		}

		playing.record.Delivery.Checks = checks
		if !announced && anyPending(checks) {
			emit(playing.sink, session.Checks, session.Data{"names": names(checks)})
			announced = true
		}

		if !anyPending(checks) {
			playing.checksEnded(checks)

			return
		}

		if time.Now().After(deadline) {
			emit(playing.sink, session.ChecksUnknown, session.Data{"reason": "still running after " + playing.checksWait.String()})

			return
		}

		select {
		case <-ctx.Done():
			emit(playing.sink, session.ChecksUnknown, session.Data{"reason": "the wait was interrupted"})

			return
		case <-time.After(playing.service.ChecksPoll):
		}
	}
}

// checksEnded prints how the checks ended.
func (playing *playing) checksEnded(checks []run.CheckResult) {
	passed, failed := 0, []string{}
	for _, check := range checks {
		switch check.Result {
		case CheckPassed:
			passed++
		case CheckFailed:
			failed = append(failed, check.Name)
		}
	}

	if len(failed) == 0 {
		emit(playing.sink, session.ChecksPassed, session.Data{"passed": passed, "total": len(checks)})

		return
	}

	emit(playing.sink, session.ChecksFailed, session.Data{"passed": passed, "total": len(checks), "failed": failed})
}

func anyPending(checks []run.CheckResult) bool {
	for _, check := range checks {
		if check.Result == CheckPending {
			return true
		}
	}

	return false
}

func names(checks []run.CheckResult) []string {
	names := make([]string, 0, len(checks))
	for _, check := range checks {
		names = append(names, check.Name)
	}

	return names
}

func (playing *playing) deliveryWords(branch Branch) string {
	words := "branch unknown"
	if branch.Name != "" {
		words = branch.Name + " at " + branch.Commit
	}
	if !branch.Upstream {
		words += ", not pushed"
	}
	if playing.record.Delivery.PullRequest != nil {
		words += ", pull request #" + strconv.Itoa(playing.record.Delivery.PullRequest.Number)
	}

	return words
}

// takeDown takes a last look at the runtime, stops it and cleans up; a
// cleanup nobody confirmed leaves the outcome unknown.
func (playing *playing) takeDown(ctx context.Context) error {
	if playing.lap.Stopped != nil {
		return nil
	}

	if last, err := playing.service.Lifecycle.Observe(ctx, playing.owner, playing.runtime); err == nil {
		playing.sample(last)
	}

	cleanup, err := playing.service.Lifecycle.Down(ctx, playing.owner, playing.runtime)
	if err != nil {
		return playing.abort(fmt.Errorf("runtime %s did not stop: %w", playing.runtime, err))
	}

	now := playing.now()
	playing.lap.Stopped = &now
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Stopped", At: now, Note: "requested by the run"})

	deployments := 0
	if playing.lap.Resources != nil {
		deployments = playing.lap.Resources.Deployments
	}

	if !cleanup.Confirmed {
		reason := fmt.Sprintf("the cleanup was not confirmed: %d of %d owned resources removed", cleanup.Removed, cleanup.Total)
		playing.lap.Outcome = Unknown
		playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Cleaned up", At: now, Note: reason})
		emit(playing.sink, session.RuntimeUnknown, session.Data{"reason": reason})

		return playing.persist()
	}

	note := fmt.Sprintf("%d of %d owned resources removed", cleanup.Removed, cleanup.Total)
	if cleanup.ClusterKept {
		note += "; cluster kept"
	}
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Cleaned up", At: now, Note: note})
	emit(playing.sink, session.RuntimeStopped, session.Data{"removed": deployments, "left": "nothing left"})

	return playing.persist()
}

// release ends the lap: the unused reservation returns to the budget, the
// run says where it stands and what comes next, and it remembers.
func (playing *playing) release() (StartResult, error) {
	now := playing.now()
	lap := &playing.lap
	if lap.Outcome != Unknown {
		lap.Outcome = Completed
	}
	lap.Ended = &now
	if gate := finops.AfterRun(true); gate != nil {
		playing.recordGate(*gate)
	}
	lap.Observed = playing.observed()

	released := 0
	if lap.Outcome == Completed {
		released = max(0, lap.Reserved-lap.ObservedCents)
		lap.Steps = append(lap.Steps, Step{Name: "Released", At: now, Note: finops.Dollars(released) + " of the reservation returned to the budget"})
	} else {
		lap.Steps = append(lap.Steps, Step{Name: "Held", At: now, Note: "the reservation is held until the outcome is known"})
	}

	outcome, proposed := playing.ending()
	playing.record.State, playing.record.Outcome = run.StateFinished, text(outcome)
	playing.record.Window.Finished = &now
	playing.record.ProposedTicketStatus = proposed

	number := 0
	if playing.record.Delivery.PullRequest != nil {
		number = playing.record.Delivery.PullRequest.Number
	}
	switch outcome {
	case ReadyForReview:
		emit(playing.sink, session.ReadyForReview, session.Data{"number": number, "releasedCents": released})
		emit(playing.sink, session.NextReview, session.Data{"number": number})
	case ReadyToMerge:
		emit(playing.sink, session.ReadyToMerge, session.Data{"number": number, "releasedCents": released})
		emit(playing.sink, session.NextMerge, nil)
	}

	playing.syncLapSummary()
	if err := playing.service.Records.WriteMemory(playing.owner, playing.runtime, playing.memory(outcome)); err != nil {
		return playing.result(), err
	}

	if err := playing.persist(); err != nil {
		return playing.result(), err
	}

	if _, err := playing.service.Index.Rebuild(playing.owner, playing.service.embedder()); err != nil {
		return playing.result(), err
	}

	return playing.result(), nil
}

// ending is how the run ends: for review when a person takes part in a
// stage the run leaves to them, or when delivery shows anything short of a
// pushed branch with a pull request whose checks all passed; to merge only
// then; unknown stays unknown.
func (playing *playing) ending() (string, *run.ProposedTicketStatus) {
	if playing.lap.Outcome == Unknown {
		return Unknown, nil
	}

	workflow := playing.guards.Workflow
	for _, stage := range workflow.Stages {
		if stage.Human == harness.OwnerOnly {
			break
		}
		if plays(stage) || !plan.IsHuman(stage, workflow.Defaults) {
			continue
		}

		playing.setStage(stage.ID, func(state *run.StageState) { state.State = run.StageWaiting })

		return ReadyForReview, &run.ProposedTicketStatus{Status: "in-review", DecidedBy: inputs.DefaultBy, Decided: false}
	}

	if !playing.deliveredCleanly() {
		return ReadyForReview, &run.ProposedTicketStatus{Status: "in-review", DecidedBy: inputs.DefaultBy, Decided: false}
	}

	return ReadyToMerge, &run.ProposedTicketStatus{Status: "ready-to-merge", DecidedBy: inputs.DefaultBy, Decided: false}
}

// deliveredCleanly reports whether delivery shows a pushed branch with a
// pull request whose every check passed.
func (playing *playing) deliveredCleanly() bool {
	delivery := playing.record.Delivery
	if !playing.pushed || delivery.PullRequest == nil || len(delivery.Checks) == 0 {
		return false
	}

	for _, check := range delivery.Checks {
		if check.Result != CheckPassed {
			return false
		}
	}

	return true
}

// observed is what the lap measured, for its fence and the memory.
func (playing *playing) observed() *run.Observed {
	observed := &run.Observed{}
	if resources := playing.lap.Resources; resources != nil {
		observed.CPUMillisPeak, observed.MemoryMiBPeak = resources.CPUPeak, resources.MemoryPeak
	}
	if !playing.lap.Unmeasured && (playing.lap.InputTokens > 0 || playing.lap.OutputTokens > 0) {
		observed.TokensThousands = &finops.TokensThousands{Input: thousands(playing.lap.InputTokens), Output: thousands(playing.lap.OutputTokens)}
	}

	return observed
}

// abort records a run that stopped in the runtime's lifecycle; the lap stays
// open and the next run plays on from the same point.
func (playing *playing) abort(cause error) error {
	playing.record.State, playing.record.Outcome = run.StateFailed, text(run.StateFailed)
	if err := playing.persist(); err != nil {
		return fmt.Errorf("%w; and the record could not be written: %v", cause, err)
	}

	return cause
}

func valueOf(value *int) int {
	if value == nil {
		return 0
	}

	return *value
}
