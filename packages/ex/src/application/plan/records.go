package plan

import (
	"errors"
	"fmt"
	"path"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// plan is the langdrift-plan fence of the runtime file.
func (planning *planning) plan() run.Plan {
	shares := run.SplitEnvelope(planning.resources, len(planning.components))

	components := make([]run.Component, 0, len(planning.components))
	for index, component := range planning.components {
		components = append(components, component.record(planning.by, shares[index]))
	}

	sources := make([]run.Source, 0, len(planning.collected.Sources))
	for _, source := range planning.collected.Sources {
		sources = append(sources, planning.sourceRecord(source))
	}

	envelope := planning.guards.Envelope
	budget := planning.guards.Budget.Mirror
	input := planning.input

	return run.Plan{
		Version: 1,
		Run:     planning.runtime.String(),
		State:   planning.state,
		Ticket: run.Ticket{
			Path:    planning.ticket.Path,
			Number:  optional(planning.ticket.Number),
			Slug:    planning.ticket.Slug,
			Title:   optional(planning.ticket.Title),
			Status:  optional(planning.ticket.Status),
			NamedBy: run.Actor{Person: planning.by, Via: inputs.Via},
		},
		Collection: run.Collection{
			Agent:   planning.collectionAgent(),
			Limits:  run.Limits{MaxDepth: planning.request.Limits.MaxDepth, MaxFiles: planning.request.Limits.MaxFiles, MaxBytes: planning.request.Limits.MaxBytes, Secrets: collection.Refused, OutsideProject: collection.Refused, Symlinks: collection.Refused},
			Reached: run.Reached{Files: planning.collected.Files, Depth: planning.collected.Depth, Bytes: planning.collected.Bytes},
			Context: run.Context{SHA256: planning.collected.ContextSHA256, Rule: collection.ContextRule},
		},
		Sources: sources,
		Facts:   nonNil(planning.analysis.facts),
		Inputs: run.Inputs{
			CPUMillis:        run.InputValue{Value: planning.resources.CPUMillis, Record: input.Record, SetBy: input.SetBy, Via: input.Via, SetAt: input.SetAt},
			MemoryMiB:        run.InputValue{Value: planning.resources.MemoryMiB, Record: input.Record, SetBy: input.SetBy, Via: input.Via, SetAt: input.SetAt},
			ModelBudgetCents: run.InputValue{Value: planning.guards.Allocation(planning.ticket), Record: budget.Record, SetBy: budget.SetBy, Via: budget.Via, SetAt: budget.SetAt},
			Ceiling:          run.Ceiling{CPUMillis: envelope.PerExecution.CPUMillis, MemoryMiB: envelope.PerExecution.MemoryMiB, Record: envelope.Mirror.Record, SetBy: envelope.Mirror.SetBy, Via: envelope.Mirror.Via, SetAt: envelope.Mirror.SetAt},
			NotRuntimeFields: nonNil(envelope.NotRuntimeFields),
		},
		Discrepancies:   nonNil(planning.analysis.discrepancies),
		Retrieval:       planning.retrieval.record(),
		Sizing:          planning.retrieval.sizing(),
		Suggestion:      planning.suggestion,
		Components:      components,
		Totals:          planning.totals,
		Materialization: run.Materialization{Deployments: len(components), ReplicasEach: 1},
		Price:           planning.price(),
		Checks:          planning.checks,
		Runtime:         run.RuntimeRef{Path: runtimePath(planning.runtime), Revision: planning.config.Revision, SHA256: planning.fenceSHA, State: planning.state},
		ContractGaps:    append([]string{}, run.ContractGaps...),
	}
}

func (planning *planning) collectionAgent() *run.CollectionAgent {
	if !planning.judge.paid() {
		return nil
	}

	return &run.CollectionAgent{Model: planning.judge.planner, PromptVersion: PromptVersion, Paid: true, LedgerLine: "planning"}
}

func (planning *planning) price() run.Price {
	price := run.Price{RuntimePriced: false, Note: "the runtime has no price source; cpu and memory stay in units"}
	if planning.guards.Prices == nil {
		return price
	}

	prices := planning.guards.Prices
	price.Record, price.Mirror, price.ReadOn = text(prices.Mirror.Record), text(".langdrift/"+planning.guards.PricesAt), text(prices.ReadOn)

	return price
}

// sourceRecord records a source by what was read and what the judge said.
func (planning *planning) sourceRecord(source collection.Source) run.Source {
	record := run.Source{
		Path:       source.Path,
		Depth:      source.Depth,
		LinkedFrom: optional(source.LinkedFrom),
		LinkSyntax: string(source.LinkSyntax),
		State:      source.State,
	}

	if source.State == collection.Read || source.State == collection.Binary {
		record.MediaType, record.Bytes, record.SHA256 = text(source.MediaType), size(source.Bytes), text(source.SHA256)
	}

	switch {
	case source.Depth == 0:
		record.InferredRole = text("ticket")
	case isGuardrail(source):
		record.InferredRole = text("guardrail")
	}

	judged, found := planning.analysis.sources[source.Path]
	if !found {
		return record
	}

	record.Verdict, record.Reason, record.Confidence = text(judged.verdict), text(judged.reason), &judged.confidence
	if judged.role != "" {
		record.InferredRole = text(judged.role)
	}

	return record
}

// runRecord is the execution record of a planned runtime: nothing has run yet.
func (planning *planning) runRecord() run.Run {
	runtime := planning.runtime.String()
	record := run.Run{
		Version:        1,
		ID:             runtime,
		Kind:           "ticket",
		Ticket:         text(planning.ticket.Path),
		Runtime:        runtime,
		State:          run.StatePlanned,
		Files:          run.Files{Plan: text("runtime.md"), Memory: "run.md", Ledger: ".langdrift/finops/ledger/" + runtime + ".ledger.json", Evidence: "evidence/"},
		HumanInTheLoop: run.HumanInTheLoop{Chosen: []string{}},
		Laps:           []run.LapSummary{},
		Stages:         []run.StageState{},
		Delivery:       run.Delivery{Checks: []run.CheckResult{}},
	}

	if planning.state == Blocked {
		record.State = Blocked
	}

	if previous := planning.previous; previous != nil {
		record.Laps, record.Delivery, record.Window = previous.Laps, previous.Delivery, run.Window{Started: previous.Window.Started}
	}

	workflow := planning.guards.Workflow
	if workflow == nil {
		return record
	}

	record.HumanInTheLoop = run.HumanInTheLoop{Record: text(workflow.Mirror.Record), Mirror: text(".langdrift/" + planning.guards.WorkflowAt), Chosen: nonNil(workflow.Defaults.HumanAt)}
	for _, stage := range workflow.Stages {
		record.Stages = append(record.Stages, run.StageState{ID: stage.ID, State: run.StagePending, Human: IsHuman(stage, workflow.Defaults)})
	}

	return record
}

// IsHuman reports whether a person takes part in a stage.
func IsHuman(stage harness.Stage, defaults harness.Defaults) bool {
	if stage.Human == harness.HumanRequired || stage.Human == harness.OwnerOnly {
		return true
	}

	for _, id := range defaults.HumanAt {
		if id == stage.ID {
			return true
		}
	}

	return false
}

// ledger is the runtime's ledger after this plan: what earlier plans and
// runs cost, with this plan's planning line added in.
func (planning *planning) ledger() finops.Ledger {
	return planning.carryOver(planning.ownLedger())
}

// ownLedger opens a ledger with this plan's planning line alone: what the
// judge's calls cost, or nothing paid when no judge ran.
func (planning *planning) ownLedger() finops.Ledger {
	budget := planning.guards.Budget
	runtime := planning.runtime.String()

	ledger := finops.Ledger{
		Mirror:   harness.Mirror{Of: inputs.Via, Record: "ledger-" + runtime, SetBy: "ex", Via: inputs.Via, SetAt: planning.now, ReadOnly: true, Note: inputs.MirrorNote},
		Version:  1,
		Run:      runtime,
		Ticket:   text(planning.ticket.Reference()),
		State:    "provisional",
		Currency: budget.Currency,
		Bounded:  finops.LedgerBounds{PerRunCents: budget.Bounded.PerRunCents, TicketAllocationCents: planning.guards.Allocation(planning.ticket), PeriodLeftCents: budget.Bounded.PeriodCents - planning.spent},
		Planning: finops.PlanningLine{Model: nil, Note: "no collection agent ran; the plan was built by the deterministic shell, nothing paid", Cents: finops.PlanningCents{Billed: number(0)}},
		Laps:     []finops.LapLine{},
		Totals:   finops.LedgerTotals{Cents: finops.TotalCents{Billed: number(0)}, PeriodLeftAfterCents: budget.Bounded.PeriodCents},
	}

	if planning.guards.Prices != nil {
		ledger.PriceSource = text(".langdrift/" + planning.guards.PricesAt)
	}
	if planning.guards.Forecasts != nil {
		ledger.ForecastSource = text(".langdrift/" + inputs.ForecastsMirror)
	}

	if !planning.judge.paid() {
		return ledger
	}

	judge := planning.judge
	observed := judge.observed()
	settlement := finops.Settle(judge.estimate, judge.reserved, observed)

	ledger.Planning = finops.PlanningLine{
		Model: text(judge.planner),
		Note:  "the planner's calls: the sources, the episodes and the suggestion; paid work under the same gates",
		Cents: finops.PlanningCents{Estimate: judge.estimate, Reserved: judge.reserved, Observed: observed, Billed: nil, Recovered: settlement.Recovered, Overspend: settlement.Overspend},
	}
	ledger.Totals = finops.LedgerTotals{
		Cents:                finops.TotalCents{Reserved: judge.reserved, Observed: observed, Billed: nil, Recovered: settlement.Recovered, Overspend: settlement.Overspend},
		PeriodLeftAfterCents: budget.Bounded.PeriodCents - observed,
	}

	return ledger
}

func (component resolved) record(by string, share harness.Resources) run.Component {
	input := component.input

	inferredFrom := []string{}
	if component.manifestName != noManifest {
		inferredFrom = append(inferredFrom, path.Join(input.Path, component.manifestName))
	}

	return run.Component{
		Path:          input.Path,
		Name:          orDefault(input.Name, componentName(input.Path)),
		Kind:          component.kind,
		Technologies:  nonNil(component.technologies),
		Image:         input.Image,
		Port:          input.Port,
		ReadinessPath: input.ReadinessPath,
		InferredFrom:  inferredFrom,
		ConfirmedBy:   run.Actor{Person: by, Via: inputs.Via},
		Proposed:      component.proposed,
		Executed:      run.Executed{CPUMillis: share.CPUMillis, MemoryMiB: share.MemoryMiB, Rule: run.EqualSplitRule},
	}
}

// loadPrevious reads the records an earlier plan of this runtime left, and
// refuses to plan over a run that is going on: a lap is open while the run
// is running, waiting or failed.
func (planning *planning) loadPrevious() error {
	record, err := planning.service.Records.ReadRun(planning.owner, planning.runtime)
	if errors.Is(err, application.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if runIsOpen(record) {
		return fmt.Errorf("%w: %s is %s with a lap open; let it end before planning again", ErrRunning, planning.runtime, record.State)
	}
	planning.previous = &record

	ledger, err := planning.service.Records.ReadLedger(planning.owner, planning.runtime.String())
	if errors.Is(err, application.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	planning.previousLedger = &ledger

	return nil
}

// runIsOpen reports whether a lap of the run is open.
func runIsOpen(record run.Run) bool {
	switch record.State {
	case run.StateRunning, run.StateWaiting, run.StateFailed:
		return true
	}

	return false
}

// loadSpent reads what the budget's period has spent so far, from the laps'
// and the plans' own state.
func (planning *planning) loadSpent() error {
	spent, err := planning.service.Spend.Period(planning.owner, planning.guards.Budget.Period)
	if err != nil {
		return err
	}
	planning.spent = spent

	return nil
}

// carryOver keeps what the earlier plans and runs of the runtime cost: the
// lap lines stay, and the planning line adds up.
func (planning *planning) carryOver(ledger finops.Ledger) finops.Ledger {
	previous := planning.previousLedger
	if previous == nil {
		return ledger
	}

	ledger.Laps = append([]finops.LapLine{}, previous.Laps...)
	ledger.Planning = addPlanning(previous.Planning, ledger.Planning)
	finops.SettleTotals(&ledger, planning.guards.Budget.Bounded.PeriodCents)

	return ledger
}

// addPlanning adds an earlier planning line to the latest.
func addPlanning(earlier, latest finops.PlanningLine) finops.PlanningLine {
	sum := latest
	sum.Cents.Estimate += earlier.Cents.Estimate
	sum.Cents.Reserved += earlier.Cents.Reserved
	sum.Cents.Observed += earlier.Cents.Observed
	sum.Cents.Recovered += earlier.Cents.Recovered
	sum.Cents.Overspend += earlier.Cents.Overspend
	sum.Cents.Billed = nil
	if earlier.Cents.Billed != nil && latest.Cents.Billed != nil {
		billed := *earlier.Cents.Billed + *latest.Cents.Billed
		sum.Cents.Billed = &billed
	}
	if latest.Model == nil {
		sum.Model = earlier.Model
	}
	if earlier.Cents.Reserved > 0 || earlier.Cents.Observed > 0 {
		sum.Note = latest.Note + "; the earlier plans of this runtime are added in"
	}

	return sum
}

// keepSpend records what this plan's calls cost: the ledger's planning
// line and the private entry the gates count. It runs once, when the plan
// is written or when a step fails after the judge was paid, so a paid call
// always leaves a money record.
func (planning *planning) keepSpend() error {
	if planning.spendKept {
		return nil
	}
	planning.spendKept = true

	if err := planning.service.Records.WriteLedger(planning.owner, planning.runtime.String(), planning.ledger()); err != nil {
		return err
	}

	own := planning.ownLedger().Planning
	if own.Cents.Reserved == 0 && own.Cents.Observed == 0 {
		return nil
	}

	return planning.service.Spend.SavePlanning(planning.owner, planning.runtime, planning.now, own)
}

// keepSpendAfter is a step's failure, with the judge's spend recorded first
// when the judge was paid.
func (planning *planning) keepSpendAfter(cause error) error {
	if !planning.judge.paid() {
		return cause
	}

	if err := planning.keepSpend(); err != nil {
		return errors.Join(cause, fmt.Errorf("and the planner's spend could not be recorded: %w", err))
	}

	return cause
}
