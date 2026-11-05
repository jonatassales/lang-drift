package plan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/collection"
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/manifest"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
	contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// planning is one plan in progress: what each act found, for the acts after it.
type planning struct {
	service Service
	request Request
	owner   project.Project
	sink    session.Sink
	guards  Guardrails
	by      string
	now     harness.Timestamp

	ticket    collection.Ticket
	collected collection.Result
	runtime   microruntime.Identifier

	connection Connection
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

	// What earlier plans and runs of the runtime left, and what the period spent.
	previous       *run.Run
	previousLedger *finops.Ledger
	spent          int
	spendKept      bool
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
// FG-01, FG-02 and FG-03. JudgeRefusal turns a missing connection or a gate
// that does not pass into a refused plan before anything is paid.
func (planning *planning) openJudge() error {
	planner := planning.guards.Planner()
	price, priced := planning.guards.Price(planner)
	estimate := planning.planningEstimate(price, priced)
	reservation := finops.RequestedReservation(estimate, planning.guards.Budget.Bounded.ReserveMarginPercent)

	if err := JudgeRefusal(planning.connection, planning.openingGate(priced, reservation)); err != nil {
		return err
	}

	planning.judge = &judge{
		client:     planning.service.Judge,
		connection: planning.connection,
		planner:    planner,
		estimate:   estimate,
		reserved:   reservation,
	}

	return nil
}

// planningEstimate is the planning forecast priced at the planner's model,
// or nothing when the forecast or the price is missing.
func (planning *planning) planningEstimate(price finops.ModelPrice, priced bool) int {
	forecasts := planning.guards.Forecasts
	if !priced || forecasts == nil || forecasts.Planning == nil {
		return 0
	}

	return finops.EstimateCents(*forecasts.Planning, price)
}

// openingGate is the verdict before the first question: a missing planning
// forecast or a gate among FG-01 to FG-03 that does not pass pauses it.
func (planning *planning) openingGate(priced bool, reservation int) AskOutcome {
	forecasts := planning.guards.Forecasts
	if forecasts == nil || forecasts.Planning == nil {
		return AskOutcome{
			Kind:   SourcesQuestion,
			Status: AskPaused,
			Reason: "no planning forecast is recorded, so nothing can be reserved; nothing was paid; record it with dft finops forecasts set",
		}
	}

	allocation := planning.guards.Allocation(planning.ticket)
	authorization := finops.Authorization{
		Priced:    priced,
		Allocated: allocation > 0,
	}
	gates := finops.BeforeAuthorization(authorization, reservation, planning.bounds())

	if !finops.Allowed(gates) {
		return AskOutcome{
			Kind:   SourcesQuestion,
			Status: AskPaused,
			Reason: gateDetail(gates) + "; nothing was paid; raise the budget",
		}
	}

	return AskOutcome{
		Kind:   SourcesQuestion,
		Status: AskCleared,
	}
}

func (planning *planning) bounds() finops.Bounds {
	budget := planning.guards.Budget.Bounded

	return finops.Bounds{PerRunCents: budget.PerRunCents, TicketAllocationCents: planning.guards.Allocation(planning.ticket), PeriodLeftCents: budget.PeriodCents - planning.spent}
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

	if err := planning.judgeSources(ctx); err != nil {
		return err
	}

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

	planning.emit(session.Budget, session.Data{"perRunCents": planning.guards.Budget.Bounded.PerRunCents, "by": planning.guards.Budget.Mirror.SetBy})

	return nil
}

// judgeSources asks the judge about every source read.
func (planning *planning) judgeSources(ctx context.Context) error {
	planning.analysis.sources = map[string]judgedSource{}

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

	question := Question{
		Kind:   SourcesQuestion,
		System: judgeSystem,
		Prompt: prompt.String(),
	}

	parsed, err := planning.judge.ask(ctx, question, func(answer contract.Answer) error {
		return checkSources(*answer.Sources, known)
	})
	if err != nil {
		return err
	}

	answer := *parsed.Sources

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

	return nil
}

// checkSources refuses an answer that judges a path it was not given, or
// gives a verdict or a confidence outside the rules.
func checkSources(answer contract.SourcesAnswer, known map[string]bool) error {
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

// envelope is the envelope the run takes: the owner's override, recorded as
// this runtime's input, or the recorded default.
func (planning *planning) envelope() (harness.Resources, run.InputValue) {
	if planning.request.CPUMillis > 0 {
		resources := harness.Resources{CPUMillis: planning.request.CPUMillis, MemoryMiB: planning.request.MemoryMiB}

		return resources, run.InputValue{Record: "runtime-" + planning.runtime.String() + "-input", SetBy: planning.by, Via: "cli", SetAt: planning.now}
	}

	mirror := planning.guards.Envelope.Mirror

	return planning.guards.Envelope.Default, run.InputValue{Record: mirror.Record, SetBy: mirror.SetBy, Via: mirror.Via, SetAt: mirror.SetAt}
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
	return finops.Grouped(cpuMillis) + " millicores and " + finops.Grouped(memoryMiB) + " MiB"
}

func figureKey(cpuMillis, memoryMiB int) string {
	return strconv.Itoa(cpuMillis) + "/" + strconv.Itoa(memoryMiB)
}

// suggest weighs the effort and the model and prints the suggestion: lines
// 22 to 24. The suggestion is the judge's; the fence keeps the recorded
// choice.
func (planning *planning) suggest(ctx context.Context) error {
	answer, err := planning.askSuggestion(ctx)
	if err != nil {
		return err
	}

	planning.emit(session.WeighEffort, session.Data{"reason": answer.EffortReason})
	planning.emit(session.WeighModel, session.Data{"reason": answer.ModelReason})

	planning.suggestion = run.Suggestion{
		By:           "judge",
		Model:        answer.Model,
		Effort:       answer.Effort,
		Reasons:      nonNil(answer.Reasons),
		AuthorizedBy: "pending",
	}

	if price, found := planning.guards.Price(planning.suggestion.Model); found {
		planning.suggestion.ModelName = text(price.Name)
	}

	planning.emit(session.Suggest, session.Data{"model": planning.suggestion.Model, "effort": planning.suggestion.Effort})

	return nil
}

// askSuggestion asks the judge for the model and the effort, from the price
// list and the efforts the harness allows.
func (planning *planning) askSuggestion(ctx context.Context) (contract.SuggestionAnswer, error) {
	if planning.guards.Prices == nil {
		return contract.SuggestionAnswer{}, errors.New("no price list is recorded to suggest a model from; set one with dft finops prices set")
	}

	models := map[string]bool{}
	var prompt strings.Builder
	prompt.WriteString("Ticket: " + planning.ticket.Title + "\n\nFacts:\n")
	for _, fact := range planning.analysis.facts {
		prompt.WriteString("- " + fact.Topic + ": " + fact.Text + "\n")
	}
	prompt.WriteString("\nComponents: " + strconv.Itoa(len(planning.request.Components)) + "\nEpisodes cited: " + strconv.Itoa(len(planning.retrieval.cited)) + "\n\nPrice list, cents per million tokens:\n")
	for _, model := range planning.guards.Prices.Models {
		models[model.ID] = true
		prompt.WriteString("- " + model.ID + ": input " + strconv.Itoa(model.InputCents) + ", output " + strconv.Itoa(model.OutputCents) + "\n")
	}
	prompt.WriteString("\nEfforts: " + strings.Join(harness.Efforts, ", ") + "\nPer-run budget: " + strconv.Itoa(planning.guards.Budget.Bounded.PerRunCents) + " cents\n\n" + suggestionPrompt)

	question := Question{
		Kind:   SuggestionQuestion,
		System: judgeSystem,
		Prompt: prompt.String(),
	}

	parsed, err := planning.judge.ask(ctx, question, func(answer contract.Answer) error {
		suggestion := answer.Suggestion
		if !models[suggestion.Model] {
			return fmt.Errorf("model %q is not on the price list", suggestion.Model)
		}
		if !harness.ValidEffort(suggestion.Effort) {
			return fmt.Errorf("effort %q is not low, medium or high", suggestion.Effort)
		}

		return nil
	})
	if err != nil {
		return contract.SuggestionAnswer{}, err
	}

	return *parsed.Suggestion, nil
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
	planning.checks = crossChecks(components, planning.resources, planning.guards.Envelope.PerExecution)
	planning.emit(session.CrossCheck, session.Data{"ports": distinctPorts(components), "images": len(components), "collisions": collisionsOf(planning.checks)})

	planning.totals = totalsOf(components, planning.resources)
	planning.emit(session.Size, session.Data{"cpuMillis": planning.totals.Proposed.CPUMillis, "boundedCpuMillis": planning.resources.CPUMillis, "memoryMiB": planning.totals.Proposed.MemoryMiB, "boundedMemoryMiB": planning.resources.MemoryMiB, "deployments": len(components)})
	planning.emit(session.Limits, session.Data{"requested": planning.limitsRequested()})
	planning.emit(session.Price, nil)

	return nil
}

// inspect reads one component's manifests, its Dockerfile and its env file,
// and sizes it. Input wins over what a file says; measured runs win over
// the sizing table.
func (planning *planning) inspect(input ComponentInput) (resolved, error) {
	inspection, err := planning.service.Inspector.Inspect(planning.owner, input.Path)
	if err != nil {
		return resolved{}, fmt.Errorf("component %s: %w", input.Path, err)
	}

	component := resolved{input: input, manifestName: orDefault(inspection.Manifest, noManifest), kind: orDefault(input.Kind, orDefault(inspection.Kind, manifest.Web))}

	if err := planning.readDockerfile(&component); err != nil {
		return resolved{}, err
	}

	if err := planning.readEnvFile(&component); err != nil {
		return resolved{}, err
	}

	input = component.input

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

// readDockerfile reads the component's Dockerfile: the port of its first
// EXPOSE when the input names none, and the readiness of its HEALTHCHECK when
// the input names no explicit kind. A missing Dockerfile is a failed check,
// not an error, so the plan still says what it found.
func (planning *planning) readDockerfile(component *resolved) error {
	input := &component.input
	component.dockerfile = input.reference().DockerfilePath()

	if input.Readiness != nil {
		probe := microruntime.ExplicitProbe(*input.Readiness)
		component.probe = &probe
	}

	content, err := planning.service.Inspector.File(planning.owner, component.dockerfile)
	if errors.Is(err, fs.ErrNotExist) {
		component.readinessErr = errors.New("the Dockerfile is missing")

		return nil
	}
	if err != nil {
		return fmt.Errorf("component %s: %w", input.Path, err)
	}

	component.dockerfileFound = true

	if input.Port == 0 {
		input.Port = manifest.ExposedPort(content)
	}

	if component.probe != nil {
		return nil
	}

	probe, err := manifest.Healthcheck(content)
	if err != nil {
		component.readinessErr = err

		return nil
	}

	component.probe = &probe

	return nil
}

// readEnvFile reads the component's env file: it must exist, hold uppercase
// names with plain-text values, and set no name the component also passes
// through. What it holds never reaches a record.
func (planning *planning) readEnvFile(component *resolved) error {
	input := component.input
	if input.EnvFile == "" {
		return nil
	}

	content, err := planning.service.Inspector.File(planning.owner, input.EnvFile)
	if errors.Is(err, fs.ErrNotExist) {
		component.envErr = errors.New("the env file is missing")

		return nil
	}
	if err != nil {
		return fmt.Errorf("component %s: %w", input.Path, err)
	}

	variables, err := configuration.ParseEnvFile(string(content))
	if err != nil {
		component.envErr = err

		return nil
	}

	component.envErr = input.reference().CheckEnv(variables)

	return nil
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
		Version:    configuration.Version,
		RuntimeUID: uid,
		RuntimeID:  planning.runtime.String(),
		ProjectID:  planning.owner.ID.String(),
		Revision:   revision,
		Name:       orDefault(planning.request.Name, orDefault(planning.ticket.Title, planning.runtime.String())),
		CPUMillis:  planning.resources.CPUMillis,
		MemoryMiB:  planning.resources.MemoryMiB,
		Model:      planning.guards.Choice.Model,
		Effort:     planning.guards.Choice.Effort,
	}

	if planning.guards.Budget.Bounded.PerRunCents > 0 {
		budget := float64(planning.guards.Budget.Bounded.PerRunCents) / 100
		config.LLMBudgetUSD = &budget
	}

	for _, component := range planning.components {
		reference := component.input.reference()
		reference.Dockerfile = reference.DockerfileName()
		config.Components = append(config.Components, reference)
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

	return planning.keepSpend()
}

func (planning *planning) emit(id string, data session.Data) {
	planning.sink.Emit(session.Event{Session: session.PlanSession, ID: id, Data: data})
}
