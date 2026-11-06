package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/inputs"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/rag"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// Service runs planned runtimes.
type Service struct {
	Projects  application.ProjectOpener
	Inputs    inputs.Store
	Records   Records
	Laps      Laps
	Driver    StageDriver
	Lifecycle Lifecycle
	Delivery  DeliveryObserver
	Index     plan.Index
	Embedder  rag.Embedder
	// ChecksPoll is how often the pull request's checks are looked at while
	// they run.
	ChecksPoll time.Duration

	// Injectable for tests; nil uses the real sources.
	IDs   application.IDSource
	Clock application.Clock
}

// StartRequest starts a lap of a planned runtime, or plays on with the lap
// that is open.
type StartRequest struct {
	Project              application.ProjectSelection
	Runtime              string
	Model                string
	Effort               string
	AcknowledgeOverspend bool
	By                   string
	Session              session.Sink
	// Confirm answers a human stage's confirmation on the spot; nil leaves
	// the run waiting for dft agent confirm.
	Confirm Confirmer
	// ChecksWait bounds how long the run follows the pull request's checks.
	ChecksWait time.Duration
}

// StartResult says where the lap stands.
type StartResult struct {
	Runtime       string `json:"runtime"`
	Lap           string `json:"lap"`
	Outcome       string `json:"outcome"`
	State         string `json:"state"`
	ReservedCents int    `json:"reservedCents"`
	ObservedCents int    `json:"observedCents"`
}

// Start opens a lap, or plays on with the open one. run.StateRunning the runtime
// authorizes its plan; the intent, the estimate and the budget print; the
// money gate reserves or pauses; then the stages play.
func (service Service) Start(ctx context.Context, request StartRequest) (StartResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return StartResult{}, err
	}

	session.Announce(request.Session, filepath.Base(owner.Root))

	runtime, err := microruntime.ParseIdentifier(request.Runtime)
	if err != nil {
		return StartResult{}, fmt.Errorf("%w: %v", inputs.ErrInvalid, err)
	}

	unlock, err := service.Laps.Lock(owner, runtime)
	if err != nil {
		return StartResult{}, fmt.Errorf("%w: %s is busy: %v", ErrRefused, runtime, err)
	}
	defer unlock()

	if !plainModel(request.Model) {
		return StartResult{}, fmt.Errorf("%w: --model takes a model id without spaces", inputs.ErrInvalid)
	}
	if !knownEffort(request.Effort) {
		return StartResult{}, fmt.Errorf("%w: --effort must be low, medium or high", inputs.ErrInvalid)
	}

	guards, err := plan.LoadGuardrails(service.Inputs, owner)
	if err != nil {
		return StartResult{}, err
	}

	config, planned, err := service.Records.ReadPlan(owner, runtime)
	if err != nil {
		return StartResult{}, err
	}

	record, err := service.Records.ReadRun(owner, runtime)
	if err != nil {
		return StartResult{}, err
	}

	if planned.State == plan.Blocked {
		return StartResult{}, fmt.Errorf("%w: the plan of %s is blocked; replan it first", ErrRefused, runtime)
	}

	if finished(record) {
		return StartResult{}, fmt.Errorf("%w: the run of %s finished, %s; plan it again for another run", ErrRefused, runtime, *record.Outcome)
	}

	ledger, err := service.Records.ReadLedger(owner, runtime.String())
	if err != nil {
		return StartResult{}, err
	}

	if guards.Workflow == nil {
		return StartResult{}, application.NotFound("no workflow is recorded; set it with dft harness workflow set")
	}

	driver, err := driverOf(service.Inputs, owner)
	if err != nil {
		return StartResult{}, err
	}
	if driver == nil && needsDriver(*guards.Workflow) {
		return StartResult{}, application.NotFound("no driver is recorded for the stages; set it with dft harness driver set")
	}

	sink := request.Session
	if sink == nil {
		sink = session.Discard{}
	}

	playing := &playing{
		service:    service,
		owner:      owner,
		runtime:    runtime,
		guards:     guards,
		driver:     driver,
		config:     config,
		planned:    planned,
		record:     record,
		ledger:     ledger,
		sink:       sink,
		by:         orDefault(request.By, inputs.DefaultBy),
		confirm:    request.Confirm,
		checksWait: request.ChecksWait,
		now:        func() harness.Timestamp { return harness.At(service.Clock.Now()) },
	}

	open, found, err := service.Laps.Open(owner, runtime)
	if err != nil {
		return StartResult{}, err
	}
	if found {
		if request.Model != "" || request.Effort != "" {
			return StartResult{}, fmt.Errorf("%w: lap %s of %s is open with %s at %s effort; a model or effort chooses a new lap once it ended", ErrRefused, open.Number, runtime, open.Choice.Model, open.Choice.Effort)
		}

		playing.lap = open
		emit(sink, session.RunScope, session.Data{"ticket": ticketReference(planned), "by": planned.Suggestion.AuthorizedBy})
		emit(sink, session.Intent, session.Data{"model": open.Choice.Model, "effort": open.Choice.Effort, "revision": config.Revision})
		emit(sink, session.LapResumed, session.Data{"lap": open.Number, "reservedCents": open.Reserved})

		return playing.play(ctx)
	}

	now := playing.now()
	by := playing.by

	if planned.State == plan.Proposed {
		planned.State, planned.Runtime.State = "authorized", "authorized"
		planned.Suggestion.AuthorizedBy = by
		planned.Authorization = &run.Authorization{By: by, Via: inputs.Via, At: now}
		if err := service.Records.WritePlan(owner, runtime, config, planned); err != nil {
			return StartResult{}, err
		}
		playing.planned = planned
	}

	lap := Lap{Number: fmt.Sprintf("%04d", len(record.Laps)+1), Run: runtime.String(), Started: now, Steps: []Step{}}
	lap.Choice = choiceOf(request, planned, by, lap.Number)
	if price, found := guards.Price(lap.Choice.Model); found {
		lap.ModelName = price.Name
	}

	emit(sink, session.RunScope, session.Data{"ticket": ticketReference(planned), "by": planned.Suggestion.AuthorizedBy})
	emit(sink, session.Intent, session.Data{"model": lap.Choice.Model, "effort": lap.Choice.Effort, "revision": config.Revision})

	forecast, forecasted := guards.Forecast(lap.Choice.Effort)
	if !forecasted {
		return StartResult{}, application.NotFound("no forecast is recorded for effort %s; set it with dft finops forecasts set", lap.Choice.Effort)
	}
	price, priced := guards.Price(lap.Choice.Model)
	lap.Estimate = forecast
	if priced {
		lap.EstimateCents = finops.EstimateCents(forecast, price)
	}
	emit(sink, session.Estimate, session.Data{"estimateCents": lap.EstimateCents, "inputThousands": forecast.Input, "outputThousands": forecast.Output})

	budget := guards.Budget
	emit(sink, session.RunBudget, session.Data{"perRunCents": budget.Bounded.PerRunCents, "by": budget.Mirror.SetBy})

	lap.Requested = finops.RequestedReservation(lap.EstimateCents, budget.Bounded.ReserveMarginPercent)
	allocation := guards.AllocationOf(ticketReference(planned))
	spent, err := service.Laps.Period(owner, budget.Period)
	if err != nil {
		return StartResult{}, err
	}
	overspent, err := service.overspent(owner, runtime)
	if err != nil {
		return StartResult{}, err
	}
	bounds := finops.Bounds{PerRunCents: budget.Bounded.PerRunCents, TicketAllocationCents: allocation, PeriodLeftCents: budget.Bounded.PeriodCents - spent}
	playing.bounds = bounds
	authorization := finops.Authorization{Priced: priced, Allocated: allocation > 0, Overspend: overspent, Acknowledged: request.AcknowledgeOverspend}
	lap.Gates = finops.BeforeAuthorization(authorization, lap.Requested, bounds)

	lap.Steps = append(lap.Steps, Step{Name: "Estimated", At: now, Note: finops.Dollars(lap.EstimateCents) + " for " + finops.Tokens(forecast.Input) + " input and " + finops.Tokens(forecast.Output) + " output tokens"})

	if !finops.Allowed(lap.Gates) {
		playing.lap, playing.record, playing.ledger = lap, record, ledger

		return playing.pause(now)
	}

	lap.Reserved = lap.Requested
	lap.Outcome = Open
	lap.Steps = append(lap.Steps, Step{Name: "Reserved", At: now, Note: finops.Dollars(lap.Reserved) + " for an estimate of " + finops.Dollars(lap.EstimateCents) + ", FG-03 passed"})
	emit(sink, session.GateOpened, session.Data{"reservedCents": lap.Reserved, "perRunCents": budget.Bounded.PerRunCents})

	if record.Window.Started == nil {
		record.Window.Started = &now
	}
	record.State, record.Outcome = run.StateRunning, text(run.StateRunning)
	record.Laps = append(record.Laps, summaryOf(lap))
	ledger.Laps = append(ledger.Laps, lineOf(lap))

	playing.lap, playing.record, playing.ledger = lap, record, ledger
	if err := playing.persist(); err != nil {
		return StartResult{}, err
	}

	return playing.play(ctx)
}

// pause records a lap that a gate stopped: nothing paid, the decision the
// owner must take printed.
func (playing *playing) pause(now harness.Timestamp) (StartResult, error) {
	lap, record, ledger, sink := &playing.lap, &playing.record, &playing.ledger, playing.sink
	bounded := playing.guards.Budget.Bounded
	gate := failing(lap.Gates)
	lap.Gate = gate.ID
	lap.Outcome = BlockedAt
	lap.Ended = &now
	ended := ErrGate
	detail := gate.ID + " " + gate.Result
	if gate.Detail != nil {
		detail = *gate.Detail
	}

	if gate.ID == "FG-03" {
		bound, boundCents := exceededBound(lap.Requested, playing.bounds)
		ended = ErrPausedAtGate
		lap.Outcome = PausedAtGate
		lap.Steps = append(lap.Steps, Step{Name: "Paused", At: now, Note: "FG-03: a " + finops.Dollars(lap.Requested) + " reservation exceeds the " + finops.Dollars(boundCents) + " " + bound})
		emit(sink, session.GatePaused, session.Data{"reservedCents": lap.Requested, "perRunCents": bounded.PerRunCents, "bound": bound, "boundCents": boundCents})
		emit(sink, session.Paused, nil)
	} else {
		lap.Steps = append(lap.Steps, Step{Name: "Blocked", At: now, Note: gate.ID + ": " + detail})
	}

	record.State, record.Outcome = run.StatePaused, text(run.StatePaused)
	record.Laps = append(record.Laps, summaryOf(*lap))
	ledger.Laps = append(ledger.Laps, lineOf(*lap))
	finops.SettleTotals(ledger, bounded.PeriodCents)

	if err := playing.service.persist(playing.owner, playing.runtime, *lap, *record, *ledger); err != nil {
		return StartResult{}, err
	}

	return StartResult{}, fmt.Errorf("%w: %s %s · %s; nothing was paid", ended, gate.ID, gate.Result, detail)
}

// persist writes the lap's private state and file, the run record and the
// ledger, in that order.
func (service Service) persist(owner project.Project, runtime microruntime.Identifier, lap Lap, record run.Run, ledger finops.Ledger) error {
	if err := service.Laps.Save(owner, runtime, lap); err != nil {
		return err
	}

	if err := service.Records.WriteLap(owner, runtime, lap); err != nil {
		return err
	}

	if err := service.Records.WriteRun(owner, runtime, record); err != nil {
		return err
	}

	return service.Records.WriteLedger(owner, runtime.String(), ledger)
}

// choiceOf is the model and effort the lap runs with: the owner's flags,
// recorded as overriding the suggestion, or the suggestion itself.
func choiceOf(request StartRequest, planned run.Plan, by, number string) run.Choice {
	choice := run.Choice{Model: planned.Suggestion.Model, Effort: planned.Suggestion.Effort, By: by, Via: inputs.Via, Record: "lap-" + number + "-choice"}

	overrode := false
	if request.Model != "" && request.Model != choice.Model {
		choice.Model, overrode = request.Model, true
	}
	if request.Effort != "" && request.Effort != choice.Effort {
		choice.Effort, overrode = request.Effort, true
	}
	if overrode {
		choice.Overrode = text("suggestion")
	}

	return choice
}

// driverOf reads the recorded driver, when the owner set one.
func driverOf(store inputs.Store, owner project.Project) (*harness.Driver, error) {
	record, err := store.Load(owner, inputs.KindDriver)
	if errors.Is(err, application.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var driver harness.Driver
	if err := json.Unmarshal(record.Value, &driver); err != nil {
		return nil, fmt.Errorf("the driver input is unreadable: %w", err)
	}

	return &driver, nil
}

// needsDriver reports whether a stage the run plays has skills to run.
func needsDriver(workflow harness.Workflow) bool {
	for _, stage := range workflow.Stages {
		if stage.Human == harness.OwnerOnly {
			return false
		}
		if len(stage.Skills) > 0 {
			return true
		}
	}

	return false
}

// ticketReference is the ticket as the ledger and the lines name it.
func ticketReference(planned run.Plan) string {
	if planned.Ticket.Number != nil {
		return *planned.Ticket.Number + "-" + planned.Ticket.Slug
	}

	return orDefault(planned.Ticket.Slug, path.Base(planned.Ticket.Path))
}

// overspent reports whether an earlier lap of the run spent more than it
// reserved, from the laps' own state: no gate reads the ledger mirror.
func (service Service) overspent(owner project.Project, runtime microruntime.Identifier) (bool, error) {
	laps, err := service.Laps.List(owner, runtime)
	if err != nil {
		return false, err
	}

	for _, lap := range laps {
		if lap.Over() && lap.ObservedCents > lap.Reserved {
			return true, nil
		}
	}

	return false, nil
}

func failing(gates []finops.GateResult) finops.GateResult {
	for _, gate := range gates {
		if gate.Result != finops.Passed {
			return gate
		}
	}

	return finops.GateResult{ID: "FG-03", Result: finops.GatePaused}
}

// summaryOf is the lap as the run record lists it.
func summaryOf(lap Lap) run.LapSummary {
	summary := run.LapSummary{Number: lap.Number, File: "laps/" + lap.Number + ".lap.md", Choice: lap.Choice, Outcome: lap.Outcome, NothingPaid: lap.NothingPaid(), EstimateCents: lap.EstimateCents, ReservedCents: lap.Reserved, ObservedCents: lap.ObservedCents}
	if lap.Gate != "" {
		summary.Gate = text(lap.Gate)
	}
	if lap.Execution != "" {
		summary.Execution = text(lap.Execution)
	}

	return summary
}

// lineOf is the lap's line of the ledger.
func lineOf(lap Lap) finops.LapLine {
	line := finops.LapLine{
		Lap:             lap.Number,
		Model:           lap.Choice.Model,
		Effort:          lap.Choice.Effort,
		TokensThousands: finops.LapTokens{Estimate: lap.Estimate},
		Cents:           finops.LapCents{Estimate: lap.EstimateCents, RequestedReservation: lap.Requested, Reserved: lap.Reserved},
		Gates:           append([]finops.GateResult{}, lap.Gates...),
		NothingPaid:     lap.NothingPaid(),
	}

	if lap.NothingPaid() {
		line.Cents.Billed = number(0)

		return line
	}

	line.Cents.Billed = nil

	return line
}

func emit(sink session.Sink, id string, data session.Data) {
	sink.Emit(session.Event{Session: session.RunSession, ID: id, Data: data})
}

func text(value string) *string { return &value }

func number(value int) *int { return &value }

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

// finished reports whether the run ended with a known outcome; a run whose
// outcome is unknown may play another lap.
func finished(record run.Run) bool {
	return record.State == run.StateFinished && record.Outcome != nil && *record.Outcome != Unknown
}

// embedder is the embedder the index is rebuilt with: the injected one, or
// the local one.
func (service Service) embedder() rag.Embedder {
	if service.Embedder != nil {
		return service.Embedder
	}

	return rag.Local{}
}

// plainModel reports whether a model flag is empty or a bare model id.
func plainModel(model string) bool {
	return model == strings.TrimSpace(model)
}

// knownEffort reports whether an effort flag is empty or one of the three efforts.
func knownEffort(effort string) bool {
	return effort == "" || harness.ValidEffort(effort)
}

// exceededBound names the first bound a reservation went past, in the order
// FG-03 checks them: the per-run cap, which the website calls the budget,
// the ticket's allocation, then what is left of the period.
func exceededBound(requested int, bounds finops.Bounds) (string, int) {
	switch {
	case requested > bounds.PerRunCents:
		return "budget", bounds.PerRunCents
	case requested > bounds.TicketAllocationCents:
		return "allocation", bounds.TicketAllocationCents
	}

	return "left of the period", bounds.PeriodLeftCents
}
