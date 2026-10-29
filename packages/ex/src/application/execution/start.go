package execution

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// Service runs planned runtimes.
type Service struct {
	Projects application.ProjectOpener
	Inputs   inputs.Store
	Records  Records
	Laps     Laps

	// Injectable for tests; nil uses the real sources.
	IDs   application.IDSource
	Clock application.Clock
}

// StartRequest starts a lap of a planned runtime.
type StartRequest struct {
	Project              application.ProjectSelection
	Runtime              string
	Model                string
	Effort               string
	AcknowledgeOverspend bool
	By                   string
	Session              session.Sink
}

// StartResult says where the lap stands.
type StartResult struct {
	Runtime       string `json:"runtime"`
	Lap           string `json:"lap"`
	Outcome       string `json:"outcome"`
	State         string `json:"state"`
	ReservedCents int    `json:"reservedCents"`
}

// The states a run may be in.
const (
	Planned = "planned"
	Running = "running"
	Paused  = "paused"
	Waiting = "waiting"
	Finished = "finished"
)

// Start opens a lap: running the runtime authorizes its plan; the intent,
// the estimate and the budget print; the money gate reserves or pauses.
func (service Service) Start(_ context.Context, request StartRequest) (StartResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return StartResult{}, err
	}

	runtime, err := microruntime.ParseIdentifier(request.Runtime)
	if err != nil {
		return StartResult{}, fmt.Errorf("%w: %v", inputs.ErrInvalid, err)
	}

	if request.Model != "" && request.Model != strings.TrimSpace(request.Model) || request.Effort != "" && !harness.ValidEffort(request.Effort) {
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

	if open, found, err := service.Laps.Open(owner, runtime); err != nil {
		return StartResult{}, err
	} else if found {
		return StartResult{}, fmt.Errorf("%w: lap %s of %s is open; stop it before another", ErrRefused, open.Number, runtime)
	}

	ledger, err := service.Records.ReadLedger(owner, runtime.String())
	if err != nil {
		return StartResult{}, err
	}

	sink := request.Session
	if sink == nil {
		sink = session.Discard{}
	}
	by := orDefault(request.By, "owner")
	now := harness.At(service.Clock.Now())

	if planned.State == plan.Proposed {
		planned.State, planned.Runtime.State = "authorized", "authorized"
		planned.Suggestion.AuthorizedBy = by
		planned.Authorization = &run.Authorization{By: by, Via: inputs.Via, At: now}
		if err := service.Records.WritePlan(owner, runtime, config, planned); err != nil {
			return StartResult{}, err
		}
	}

	lap := Lap{Number: fmt.Sprintf("%04d", len(record.Laps)+1), Run: runtime.String(), Started: now, Steps: []Step{}}
	lap.Choice = choiceOf(request, planned, by, lap.Number)
	if price, found := guards.Price(lap.Choice.Model); found {
		lap.ModelName = price.Name
	}

	emit(sink, session.RunCommand, session.Data{"runtime": runtime.String()})
	emit(sink, session.RunScope, session.Data{"ticket": ticketReference(planned), "by": planned.Suggestion.AuthorizedBy})
	emit(sink, session.Intent, session.Data{"model": lap.Choice.Model, "effort": lap.Choice.Effort, "revision": config.Revision})

	forecast, forecasted := guards.Forecast(lap.Choice.Effort)
	if !forecasted {
		return StartResult{}, application.NotFound("no forecast is recorded for effort %s; set it with ex finops forecasts set", lap.Choice.Effort)
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
	bounds := finops.Bounds{PerRunCents: budget.Bounded.PerRunCents, TicketAllocationCents: allocation, PeriodLeftCents: budget.Bounded.PeriodCents - ledger.Totals.Cents.Observed}
	authorization := finops.Authorization{Priced: priced, Allocated: allocation > 0, Overspend: overspent(ledger), Acknowledged: request.AcknowledgeOverspend}
	lap.Gates = finops.BeforeAuthorization(authorization, lap.Requested, bounds)

	lap.Steps = append(lap.Steps, Step{Name: "Estimated", At: now, Note: finops.Dollars(lap.EstimateCents) + " for " + finops.Tokens(forecast.Input) + " input and " + finops.Tokens(forecast.Output) + " output tokens"})

	if !finops.Allowed(lap.Gates) {
		return service.pause(owner, runtime, &lap, &record, &ledger, sink, now, budget.Bounded.PerRunCents)
	}

	lap.Reserved = lap.Requested
	lap.Outcome = Open
	lap.Steps = append(lap.Steps, Step{Name: "Reserved", At: now, Note: finops.Dollars(lap.Reserved) + " for an estimate of " + finops.Dollars(lap.EstimateCents) + ", FG-03 passed"})
	emit(sink, session.GateOpened, session.Data{"reservedCents": lap.Reserved, "perRunCents": budget.Bounded.PerRunCents})

	if record.Window.Started == nil {
		record.Window.Started = &now
	}
	record.State, record.Outcome = Running, text("running")
	record.Laps = append(record.Laps, summaryOf(lap))
	ledger.Laps = append(ledger.Laps, lineOf(lap))
	settleTotals(&ledger, budget.Bounded.PeriodCents)

	if err := service.persist(owner, runtime, lap, record, ledger); err != nil {
		return StartResult{}, err
	}

	return StartResult{Runtime: runtime.String(), Lap: lap.Number, Outcome: Open, State: record.State, ReservedCents: lap.Reserved}, nil
}

// pause records a lap that a gate stopped: nothing paid, the decision the
// owner must take printed.
func (service Service) pause(owner project.Project, runtime microruntime.Identifier, lap *Lap, record *run.Run, ledger *finops.Ledger, sink session.Sink, now harness.Timestamp, perRunCents int) (StartResult, error) {
	gate := failing(lap.Gates)
	lap.Gate = gate.ID
	lap.Outcome = BlockedAt
	lap.Ended = &now
	detail := gate.ID + " " + gate.Result
	if gate.Detail != nil {
		detail = *gate.Detail
	}

	if gate.ID == "FG-03" {
		lap.Outcome = PausedAtGate
		lap.Steps = append(lap.Steps, Step{Name: "Paused", At: now, Note: "FG-03: a " + finops.Dollars(lap.Requested) + " reservation exceeds the " + finops.Dollars(perRunCents) + " per-run cap"})
		emit(sink, session.GatePaused, session.Data{"reservedCents": lap.Requested, "perRunCents": perRunCents})
		emit(sink, session.Paused, nil)
	} else {
		lap.Steps = append(lap.Steps, Step{Name: "Blocked", At: now, Note: gate.ID + ": " + detail})
		emit(sink, session.GateBlocked, session.Data{"gate": gate.ID, "detail": detail})
	}

	record.State, record.Outcome = Paused, text(Paused)
	record.Laps = append(record.Laps, summaryOf(*lap))
	ledger.Laps = append(ledger.Laps, lineOf(*lap))
	settleTotals(ledger, ledger.Bounded.PeriodLeftCents)

	if err := service.persist(owner, runtime, *lap, *record, *ledger); err != nil {
		return StartResult{}, err
	}

	return StartResult{}, fmt.Errorf("%w: %s %s · %s; nothing was paid", ErrGate, gate.ID, gate.Result, detail)
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

// ticketReference is the ticket as the ledger and the lines name it.
func ticketReference(planned run.Plan) string {
	if planned.Ticket.Number != nil {
		return *planned.Ticket.Number + "-" + planned.Ticket.Slug
	}

	return orDefault(planned.Ticket.Slug, path.Base(planned.Ticket.Path))
}

// overspent reports whether an earlier lap of the run overspent.
func overspent(ledger finops.Ledger) bool {
	for _, line := range ledger.Laps {
		if line.Cents.Overspend > 0 {
			return true
		}
	}

	return false
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
	summary := run.LapSummary{Number: lap.Number, File: "laps/" + lap.Number + ".lap.md", Choice: lap.Choice, Outcome: lap.Outcome, NothingPaid: lap.NothingPaid(), EstimateCents: lap.EstimateCents, ReservedCents: lap.Reserved}
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

// settleTotals sums the ledger's lines. Billed stays unknown while any line's is.
func settleTotals(ledger *finops.Ledger, periodCents int) {
	totals := finops.TotalCents{Reserved: ledger.Planning.Cents.Reserved, Observed: ledger.Planning.Cents.Observed, Recovered: ledger.Planning.Cents.Recovered, Overspend: ledger.Planning.Cents.Overspend}
	billed, known := 0, ledger.Planning.Cents.Billed != nil
	if known {
		billed = *ledger.Planning.Cents.Billed
	}

	for _, line := range ledger.Laps {
		totals.Reserved += line.Cents.Reserved
		totals.Observed += line.Cents.Observed
		totals.Recovered += line.Cents.Recovered
		totals.Overspend += line.Cents.Overspend
		if line.Cents.Billed == nil {
			known = false

			continue
		}
		billed += *line.Cents.Billed
	}

	if known {
		totals.Billed = number(billed)
	}

	ledger.Totals = finops.LedgerTotals{Cents: totals, PeriodLeftAfterCents: periodCents - totals.Observed}
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
