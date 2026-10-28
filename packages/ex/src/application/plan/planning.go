package plan

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// planning is one plan in progress: what each act found, for the acts after it.
type planning struct {
	service Service
	request Request
	owner   project.Project
	sink    session.Sink
	guards  guardrails
	by      string
	now     harness.Timestamp

	ticket    collection.Ticket
	collected collection.Result
	runtime   microruntime.Identifier

	judge      *judge
	analysis   analysis
	retrieval  retrieved
	suggestion run.Suggestion

	resources  harness.Resources
	input      run.InputValue
	components []resolved
	totals     run.Totals
	checks     []run.Check
	state      string
	config     configuration.Configuration
	fenceSHA   string
}

// analysis is what the judge said about the sources.
type analysis struct {
	sources       map[string]judgedSource
	facts         []run.Fact
	criteria      []string
	figures       []figure
	discrepancies []run.Discrepancy
}

// figure is a CPU and memory pair the material states, with its citation.
type figure struct {
	topic  string
	fact   string
	source string
	key    string
}

// judgedSource is the judge's verdict on one source.
type judgedSource struct {
	verdict    string
	reason     string
	confidence float64
	role       string
	summary    string
}

// collect reads the ticket and follows its links: lines 1 to 4.
func (planning *planning) collect(ticketPath string) error {
	planning.emit(session.PlanCommand, session.Data{"ticket": path.Base(ticketPath)})
	planning.emit(session.PlanScope, nil)
	planning.emit(session.ReadTicket, session.Data{"ticket": path.Base(ticketPath)})

	collected, err := planning.service.Collector.Collect(planning.owner, ticketPath, planning.request.Limits)
	if err != nil {
		return err
	}

	if len(collected.Sources) == 0 || collected.Sources[0].State != collection.Read {
		return application.NotFound("no ticket could be read at %s", ticketPath)
	}

	planning.collected = collected
	planning.ticket = collection.ParseTicket(ticketPath, collected.Sources[0].Content)

	runtime, err := runtimeIDOf(planning.request.Runtime, planning.ticket.Slug)
	if err != nil {
		return err
	}
	planning.runtime = runtime

	planning.emit(session.Links, session.Data{"status": orDefault(planning.ticket.Status, "no status"), "links": linkNames(collected)})

	return nil
}

// openJudge decides whether the planner's paid calls may run: the planning
// forecast priced at the planner's model, reserved with the margin, through
// FG-01, FG-02 and FG-03. A gate that does not pass refuses the plan before
// anything is paid; a judge that cannot run leaves the deterministic shell.
func (planning *planning) openJudge() error {
	reason := planning.judgeUnavailable()
	if reason != "" {
		planning.emit(session.JudgeUnavailable, session.Data{"reason": reason})

		return nil
	}

	planner := planning.guards.planner()
	price, priced := planning.guards.price(planner)
	estimate := 0
	if priced {
		estimate = finops.EstimateCents(*planning.guards.forecasts.Planning, price)
	}
	reservation := finops.RequestedReservation(estimate, planning.guards.budget.Bounded.ReserveMarginPercent)

	allocation := planning.guards.allocation(planning.ticket)
	gates := finops.BeforeAuthorization(finops.Authorization{Priced: priced, Allocated: allocation > 0}, reservation, planning.bounds())
	if !finops.Allowed(gates) {
		return fmt.Errorf("%w: %s; nothing was paid; raise the budget or plan with --no-judge", ErrPaused, gateDetail(gates))
	}

	planning.judge = &judge{model: planning.service.Model, planner: planner, price: price, estimate: estimate, reserved: reservation}

	return nil
}

// judgeUnavailable says why the judge cannot run, or nothing when it can.
func (planning *planning) judgeUnavailable() string {
	switch {
	case planning.request.NoJudge:
		return "--no-judge; verdicts, facts and the suggestion are the inputs' own"
	case planning.service.Model == nil:
		return "no model is configured; verdicts, facts and the suggestion are the inputs' own"
	case planning.guards.forecasts == nil || planning.guards.forecasts.Planning == nil:
		return "no planning forecast is recorded; set planning with ex finops forecasts set"
	}

	return ""
}

func (planning *planning) bounds() finops.Bounds {
	budget := planning.guards.budget.Bounded

	return finops.Bounds{PerRunCents: budget.PerRunCents, TicketAllocationCents: planning.guards.allocation(planning.ticket), PeriodLeftCents: budget.PeriodCents}
}

func gateDetail(gates []finops.GateResult) string {
	for _, gate := range gates {
		if gate.Result == finops.Passed {
			continue
		}
		if gate.Detail != nil {
			return gate.ID + " " + gate.Result + " · " + *gate.Detail
		}

		return gate.ID + " " + gate.Result
	}

	return "a gate did not pass"
}

// analyze judges the sources and prints the analysis act: lines 5 to 13.
func (planning *planning) analyze(ctx context.Context) error {
	planning.emit(session.Analysis, nil)
	planning.judgeSources(ctx)

	for _, source := range planning.collected.Sources[1:] {
		if source.State != collection.Read && source.State != collection.Binary {
			continue
		}

		summary := source.State
		if judged, found := planning.analysis.sources[source.Path]; found && judged.summary != "" {
			summary = judged.summary
		}

		planning.emit(session.ReadDocument, session.Data{"name": path.Base(source.Path), "path": source.Path, "summary": summary})
	}

	planning.emit(session.Reads, session.Data{"count": planning.collected.Files, "refused": len(planning.collected.Sources) - planning.collected.Files})

	for _, source := range planning.collected.Sources {
		if isGuardrail(source) {
			planning.emit(session.Guardrail, session.Data{"path": source.Path, "criteria": criteria(source.Content)})
		}
	}

	if len(planning.analysis.criteria) > 0 {
		planning.emit(session.JudgeGuardrail, session.Data{"judgements": planning.analysis.criteria})
	}

	planning.resources, planning.input = planning.envelope()
	planning.emit(session.Bounds, session.Data{"cpuMillis": planning.resources.CPUMillis, "memoryMiB": planning.resources.MemoryMiB, "by": planning.input.SetBy, "record": planning.input.Record})

	for _, differing := range planning.compare() {
		planning.emit(session.Discrepancy, session.Data{"source": differing.source, "fact": differing.fact})
	}

	planning.emit(session.Budget, session.Data{"perRunCents": planning.guards.budget.Bounded.PerRunCents, "by": planning.guards.budget.Mirror.SetBy})

	return nil
}

// judgeSources asks the judge about every source read.
func (planning *planning) judgeSources(ctx context.Context) {
	planning.analysis.sources = map[string]judgedSource{}
	if planning.judge == nil {
		return
	}

	var prompt strings.Builder
	prompt.WriteString("Ticket: " + planning.ticket.Path + "\n\n")

	known := map[string]bool{}
	for _, source := range planning.collected.Sources {
		if source.State != collection.Read {
			continue
		}
		known[source.Path] = true
		prompt.WriteString(material(source.Path, string(source.Content), 12000) + "\n\n")
	}
	prompt.WriteString(sourcesPrompt)

	var answer sourcesAnswer
	accepted := planning.judge.ask(ctx, Question{Kind: SourcesQuestion, System: judgeSystem, Prompt: prompt.String()}, &answer, func() error {
		return checkSources(answer, known)
	})
	planning.stopIfJudgeStopped()

	if !accepted {
		return
	}

	for _, source := range answer.Sources {
		planning.analysis.sources[source.Path] = judgedSource{verdict: source.Verdict, reason: source.Reason, confidence: source.Confidence, role: source.Role, summary: source.Summary}
	}

	for _, fact := range answer.Facts {
		planning.analysis.facts = append(planning.analysis.facts, run.Fact{Topic: fact.Topic, Text: fact.Text, Source: fact.Source})
	}

	for _, criterion := range answer.Criteria {
		planning.analysis.criteria = append(planning.analysis.criteria, criterion.ID+" "+criterion.Applies)
	}

	for _, stated := range answer.Figures {
		planning.analysis.figures = append(planning.analysis.figures, figure{topic: stated.Topic, fact: stated.Fact, source: stated.Source, key: figureKey(stated.CPUMillis, stated.MemoryMiB)})
	}
}

// checkSources refuses an answer that judges a path it was not given, or
// gives a verdict or a confidence outside the rules.
func checkSources(answer sourcesAnswer, known map[string]bool) error {
	for _, source := range answer.Sources {
		if !known[source.Path] {
			return fmt.Errorf("it judges %q, which was not read", source.Path)
		}
		if source.Verdict != "kept" && source.Verdict != "discarded" {
			return fmt.Errorf("verdict %q is neither kept nor discarded", source.Verdict)
		}
		if source.Confidence < 0 || source.Confidence > 1 {
			return fmt.Errorf("confidence %v is outside 0 to 1", source.Confidence)
		}
	}

	for _, fact := range answer.Facts {
		cited, _, _ := strings.Cut(fact.Source, "#")
		if !known[cited] {
			return fmt.Errorf("a fact cites %q, which was not read", fact.Source)
		}
	}

	return nil
}

// stopIfJudgeStopped prints, once, why the judge stopped.
func (planning *planning) stopIfJudgeStopped() {
	if planning.judge == nil || planning.judge.stopped == "" || planning.judge.stopped == "printed" {
		return
	}

	planning.emit(session.JudgeStopped, session.Data{"reason": planning.judge.stopped})
	planning.judge.stopped = "printed"
}

// envelope is the envelope the run takes: the owner's override, recorded as
// this runtime's input, or the recorded default.
func (planning *planning) envelope() (harness.Resources, run.InputValue) {
	if planning.request.CPUMillis > 0 {
		resources := harness.Resources{CPUMillis: planning.request.CPUMillis, MemoryMiB: planning.request.MemoryMiB}

		return resources, run.InputValue{Record: "runtime-" + planning.runtime.String() + "-input", SetBy: planning.by, Via: "cli", SetAt: planning.now}
	}

	mirror := planning.guards.envelope.Mirror

	return planning.guards.envelope.Default, run.InputValue{Record: mirror.Record, SetBy: mirror.SetBy, Via: mirror.Via, SetAt: mirror.SetAt}
}

// compare sets each figure the material gave against the input; the input
// holds either way. It returns the figures that differ.
func (planning *planning) compare() []figure {
	input := figureText(planning.resources.CPUMillis, planning.resources.MemoryMiB)
	inputKey := figureKey(planning.resources.CPUMillis, planning.resources.MemoryMiB)

	var differing []figure
	for _, stated := range planning.analysis.figures {
		result := "consistent"
		if stated.key != inputKey {
			result = "differs"
			differing = append(differing, stated)
		}

		planning.analysis.discrepancies = append(planning.analysis.discrepancies, run.Discrepancy{Topic: stated.topic, Fact: stated.fact, Input: input, Result: result})
	}

	return differing
}

func figureText(cpuMillis, memoryMiB int) string {
	return grouped(cpuMillis) + " millicores and " + grouped(memoryMiB) + " MiB"
}

func figureKey(cpuMillis, memoryMiB int) string {
	return strconv.Itoa(cpuMillis) + "/" + strconv.Itoa(memoryMiB)
}

// suggest weighs the effort and the model and prints the suggestion: lines
// 22 to 24. Without a judge, the suggestion is the recorded choice.
func (planning *planning) suggest(ctx context.Context) {
	planning.suggestion = run.Suggestion{By: "input", Model: planning.guards.choice.Model, Effort: planning.guards.choice.Effort, Reasons: []string{}, AuthorizedBy: "pending"}

	if answer, accepted := planning.askSuggestion(ctx); accepted {
		planning.emit(session.WeighEffort, session.Data{"reason": answer.EffortReason})
		planning.emit(session.WeighModel, session.Data{"reason": answer.ModelReason})
		planning.suggestion.By = "recommendation"
		planning.suggestion.Model = answer.Model
		planning.suggestion.Effort = answer.Effort
		planning.suggestion.Reasons = nonNil(answer.Reasons)
	}

	if price, found := planning.guards.price(planning.suggestion.Model); found {
		planning.suggestion.ModelName = text(price.Name)
	}

	planning.emit(session.Suggest, session.Data{"model": planning.suggestion.Model, "effort": planning.suggestion.Effort})
}

func (planning *planning) askSuggestion(ctx context.Context) (suggestionAnswer, bool) {
	var answer suggestionAnswer
	if planning.judge == nil || planning.guards.prices == nil {
		return answer, false
	}

	models := map[string]bool{}
	var prompt strings.Builder
	prompt.WriteString("Ticket: " + planning.ticket.Title + "\n\nFacts:\n")
	for _, fact := range planning.analysis.facts {
		prompt.WriteString("- " + fact.Topic + ": " + fact.Text + "\n")
	}
	prompt.WriteString("\nComponents: " + strconv.Itoa(len(planning.request.Components)) + "\nEpisodes cited: " + strconv.Itoa(len(planning.retrieval.cited)) + "\n\nPrice list, cents per million tokens:\n")
	for _, model := range planning.guards.prices.Models {
		models[model.ID] = true
		prompt.WriteString("- " + model.ID + ": input " + strconv.Itoa(model.InputCents) + ", output " + strconv.Itoa(model.OutputCents) + "\n")
	}
	prompt.WriteString("\nEfforts: " + strings.Join(harness.Efforts, ", ") + "\nPer-run budget: " + strconv.Itoa(planning.guards.budget.Bounded.PerRunCents) + " cents\n\n" + suggestionPrompt)

	accepted := planning.judge.ask(ctx, Question{Kind: SuggestionQuestion, System: judgeSystem, Prompt: prompt.String()}, &answer, func() error {
		if !models[answer.Model] {
			return fmt.Errorf("model %q is not on the price list", answer.Model)
		}
		if !harness.ValidEffort(answer.Effort) {
			return fmt.Errorf("effort %q is not low, medium or high", answer.Effort)
		}

		return nil
	})
	planning.stopIfJudgeStopped()

	return answer, accepted
}

// resolve inspects and sizes the components: lines 25 to 38.
func (planning *planning) resolve() error {
	planning.emit(session.Planning, nil)
	planning.emit(session.Resolve, session.Data{"count": len(planning.request.Components), "namedIn": "the CLI"})

	components := make([]resolved, 0, len(planning.request.Components))
	for _, input := range planning.request.Components {
		component, err := planning.inspect(input)
		if err != nil {
			return err
		}
		components = append(components, component)
	}

	planning.components = components
	planning.checks = crossChecks(components, planning.resources, planning.guards.envelope.PerExecution)
	planning.emit(session.CrossCheck, session.Data{"ports": distinctPorts(components), "images": len(components), "collisions": collisionsOf(planning.checks)})

	planning.totals = totalsOf(components, planning.resources)
	planning.emit(session.Size, session.Data{"cpuMillis": planning.totals.Proposed.CPUMillis, "boundedCpuMillis": planning.resources.CPUMillis, "memoryMiB": planning.totals.Proposed.MemoryMiB, "boundedMemoryMiB": planning.resources.MemoryMiB, "deployments": len(components)})
	planning.emit(session.Limits, session.Data{"requested": planning.limitsRequested()})
	planning.emit(session.Price, nil)

	return nil
}

// inspect reads one component's manifests and sizes it. Input wins over
// what a manifest says; measured runs win over the sizing table.
func (planning *planning) inspect(input ComponentInput) (resolved, error) {
	inspection, err := planning.service.Inspector.Inspect(planning.owner, input.Path)
	if err != nil {
		return resolved{}, fmt.Errorf("component %s: %w", input.Path, err)
	}

	component := resolved{input: input, manifestName: orDefault(inspection.Manifest, noManifest), kind: orDefault(input.Kind, orDefault(inspection.Kind, manifest.Web))}

	technologies := inspection.Technologies
	if len(input.Technologies) > 0 {
		technologies = nil
		for _, technology := range input.Technologies {
			technologies = append(technologies, manifest.Technology{ID: technology.ID, Label: technology.Label})
		}
	}
	for _, technology := range technologies {
		component.technologies = append(component.technologies, run.Technology{ID: technology.ID, Label: technology.Label})
	}

	component.proposed = manifest.Proposed(component.kind, technologies)
	if planning.retrieval.measured != nil {
		component.proposed = planning.retrieval.share(len(planning.request.Components))
	}

	planning.emit(session.Inspect, session.Data{"path": input.Path, "manifest": component.manifestName, "port": input.Port})
	planning.emit(session.Component, session.Data{"path": input.Path, "technology": manifest.Labels(technologies), "kind": component.kind, "cpuMillis": component.proposed.CPUMillis, "memoryMiB": component.proposed.MemoryMiB})

	return component, nil
}

// limitsRequested says whether the material asked for a resource runtime.md
// cannot express.
func (planning *planning) limitsRequested() string {
	for _, fact := range planning.analysis.facts {
		lower := strings.ToLower(fact.Text)
		if strings.Contains(lower, "gpu") || strings.Contains(lower, "disk") {
			return "requested in " + fact.Source
		}
	}

	return "none requested"
}

// check builds the fence, validates it and decides the plan's state: lines
// 39 to 43.
func (planning *planning) check() error {
	revision, uid, err := planning.revision()
	if err != nil {
		return err
	}

	planning.config = planning.configuration(uid, revision)

	schema := run.Check{ID: "schema", Result: passed, Errors: number(0)}
	if err := planning.config.Validate(); err != nil {
		schema = run.Check{ID: "schema", Result: failed, Errors: number(1), Detail: text(err.Error())}
	}
	planning.checks = append([]run.Check{schema}, planning.checks...)
	planning.emit(session.Validate, session.Data{"errors": *schema.Errors})

	planning.state = Proposed
	if failing := failingChecks(planning.checks); len(failing) > 0 {
		planning.state = Blocked
		planning.emit(session.PlanBlocked, session.Data{"revision": revision, "checks": failing})
	} else {
		planning.emit(session.Draft, session.Data{"revision": revision})
		planning.emit(session.Handoff, nil)
	}

	_, fenceSHA, err := planning.service.Records.Fence(planning.config)
	planning.fenceSHA = fenceSHA

	return err
}

// revision is the next revision of the runtime, and its uid, kept across
// revisions. A fence that no longer validates still names its revision.
func (planning *planning) revision() (int, string, error) {
	existing, err := planning.service.Records.ReadRuntime(planning.owner, planning.runtime)
	if errors.Is(err, application.ErrNotFound) {
		uid, err := planning.service.IDs.ID()

		return 1, uid, err
	}
	if err != nil {
		return 0, "", err
	}

	config, err := connection.Parse(string(existing))
	if config.Revision < 1 {
		return 0, "", fmt.Errorf("the existing runtime.md of %s does not parse: %w", planning.runtime, err)
	}

	return config.Revision + 1, config.RuntimeUID, nil
}

// configuration is the fence the engine executes. The model and effort are
// the recorded choice: a suggestion never sets one.
func (planning *planning) configuration(uid string, revision int) configuration.Configuration {
	config := configuration.Configuration{
		Version:    1,
		RuntimeUID: uid,
		RuntimeID:  planning.runtime.String(),
		ProjectID:  planning.owner.ID.String(),
		Revision:   revision,
		Name:       orDefault(planning.request.Name, orDefault(planning.ticket.Title, planning.runtime.String())),
		CPUMillis:  planning.resources.CPUMillis,
		MemoryMiB:  planning.resources.MemoryMiB,
		Model:      planning.guards.choice.Model,
		Effort:     planning.guards.choice.Effort,
	}

	if planning.guards.budget.Bounded.PerRunCents > 0 {
		budget := float64(planning.guards.budget.Bounded.PerRunCents) / 100
		config.LLMBudgetUSD = &budget
	}

	for _, component := range planning.request.Components {
		build := component.Build
		if build.Args == nil {
			build.Args = []string{}
		}
		config.Components = append(config.Components, configuration.Component{Path: component.Path, Image: component.Image, Port: component.Port, ReadinessPath: component.ReadinessPath, Build: &build})
	}

	return config
}

// write renders and writes the runtime file, the run record and the ledger.
func (planning *planning) write() error {
	document := Document{
		Name:          planning.config.Name,
		TicketName:    planning.ticket.Name(),
		Run:           planning.runtime.String(),
		ContextSHA256: planning.collected.ContextSHA256,
		Record:        planning.input.Record,
		By:            planning.input.SetBy,
		Date:          planning.input.SetAt.Format("2006-01-02"),
		Config:        planning.config,
		Plan:          planning.plan(),
	}

	markdown, err := planning.service.Records.Render(document)
	if err != nil {
		return err
	}

	if err := planning.service.Records.WriteRuntime(planning.owner, planning.runtime, markdown); err != nil {
		return err
	}

	if err := planning.service.Records.WriteRun(planning.owner, planning.runtime, planning.runRecord()); err != nil {
		return err
	}

	return planning.service.Records.WriteLedger(planning.owner, planning.runtime.String(), planning.ledger())
}

func (planning *planning) emit(id string, data session.Data) {
	planning.sink.Emit(session.Event{Session: session.PlanSession, ID: id, Data: data})
}
