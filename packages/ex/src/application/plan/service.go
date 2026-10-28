package plan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
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

// Service plans tickets.
type Service struct {
	Projects  application.ProjectOpener
	Collector Collector
	Inspector Inspector
	Inputs    inputs.Store
	Records   Records

	// Injectable for tests; nil uses the real sources.
	IDs   application.IDSource
	Clock application.Clock
}

// ComponentInput is one component the owner names on the command line.
type ComponentInput struct {
	Path          string
	Name          string
	Kind          string
	Image         string
	Port          int
	ReadinessPath string
	Build         configuration.Build
	Technologies  []run.Technology
}

// Request plans one ticket.
type Request struct {
	Project    application.ProjectSelection
	Ticket     string
	Runtime    string
	Name       string
	Components []ComponentInput
	CPUMillis  int
	MemoryMiB  int
	Limits     collection.Limits
	NoJudge    bool
	By         string
	Session    session.Sink
}

// Result says what was written.
type Result struct {
	Runtime  string `json:"runtime"`
	Revision int    `json:"revision"`
	Path     string `json:"path"`
	State    string `json:"state"`
}

// The states a plan may be in.
const (
	Proposed = "proposed"
	Blocked  = "blocked"
)

// The results a check may have.
const (
	passed = "passed"
	failed = "failed"
)

var criterionID = regexp.MustCompile(`\bFIN-\d{2}\b`)

// Plan reads the ticket, follows its links, takes the inputs, resolves and
// sizes the components, checks the draft and writes the runtime file, the
// run record and the ledger's planning line.
func (service Service) Plan(_ context.Context, request Request) (Result, error) {
	if err := validate(request); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	guards, err := service.guardrails(owner)
	if err != nil {
		return Result{}, err
	}

	ticketPath, err := ticketPathOf(owner, request.Ticket)
	if err != nil {
		return Result{}, err
	}

	sink := request.Session
	if sink == nil {
		sink = session.Discard{}
	}
	by := orDefault(request.By, "owner")
	now := harness.At(service.Clock.Now())

	emit(sink, session.PlanCommand, session.Data{"ticket": path.Base(ticketPath)})
	emit(sink, session.PlanScope, nil)
	emit(sink, session.ReadTicket, session.Data{"ticket": path.Base(ticketPath)})

	collected, err := service.Collector.Collect(owner, ticketPath, request.Limits)
	if err != nil {
		return Result{}, err
	}

	if len(collected.Sources) == 0 || collected.Sources[0].State != collection.Read {
		return Result{}, application.NotFound("no ticket could be read at %s", ticketPath)
	}

	ticket := collection.ParseTicket(ticketPath, collected.Sources[0].Content)

	runtimeID, err := runtimeIDOf(request.Runtime, ticket.Slug)
	if err != nil {
		return Result{}, err
	}

	emit(sink, session.Links, session.Data{"status": orDefault(ticket.Status, "no status"), "links": linkNames(collected)})
	emit(sink, session.Analysis, nil)

	for _, source := range collected.Sources[1:] {
		if source.State != collection.Read && source.State != collection.Binary {
			continue
		}
		emit(sink, session.ReadDocument, session.Data{"name": path.Base(source.Path), "path": source.Path, "summary": source.State})
	}

	emit(sink, session.Reads, session.Data{"count": collected.Files, "refused": len(collected.Sources) - collected.Files})

	for _, source := range collected.Sources {
		if !isGuardrail(source) {
			continue
		}
		emit(sink, session.Guardrail, session.Data{"path": source.Path, "criteria": criteria(source.Content)})
	}

	emit(sink, session.JudgeUnavailable, session.Data{"reason": judgeReason(request)})

	resources, record := bounds(request, guards, runtimeID, by, now)
	emit(sink, session.Bounds, session.Data{"cpuMillis": resources.CPUMillis, "memoryMiB": resources.MemoryMiB, "by": record.SetBy, "record": record.Record})
	emit(sink, session.Budget, session.Data{"perRunCents": guards.budget.Bounded.PerRunCents, "by": guards.budget.Mirror.SetBy})

	emit(sink, session.Planning, nil)
	emit(sink, session.Resolve, session.Data{"count": len(request.Components), "namedIn": "the CLI"})

	components, err := service.resolve(owner, request.Components, sink)
	if err != nil {
		return Result{}, err
	}

	checks := crossChecks(components, resources, guards.envelope.PerExecution)
	emit(sink, session.CrossCheck, session.Data{"ports": distinctPorts(components), "images": len(components), "collisions": collisionsOf(checks)})

	totals := totalsOf(components, resources)
	emit(sink, session.Size, session.Data{"cpuMillis": totals.Proposed.CPUMillis, "boundedCpuMillis": resources.CPUMillis, "memoryMiB": totals.Proposed.MemoryMiB, "boundedMemoryMiB": resources.MemoryMiB, "deployments": len(components)})
	emit(sink, session.Limits, session.Data{"requested": "none requested"})
	emit(sink, session.Price, nil)

	revision, runtimeUID, err := service.revision(owner, runtimeID)
	if err != nil {
		return Result{}, err
	}

	config := service.configuration(request, guards, owner, runtimeID, runtimeUID, revision, resources, orDefault(request.Name, orDefault(ticket.Title, runtimeID.String())))

	schema := run.Check{ID: "schema", Result: passed, Errors: number(0)}
	if err := config.Validate(); err != nil {
		schema = run.Check{ID: "schema", Result: failed, Errors: number(1), Detail: text(err.Error())}
	}
	checks = append([]run.Check{schema}, checks...)
	emit(sink, session.Validate, session.Data{"errors": *schema.Errors})

	state, failing := Proposed, failingChecks(checks)
	if len(failing) > 0 {
		state = Blocked
		emit(sink, session.PlanBlocked, session.Data{"revision": revision, "checks": failing})
	} else {
		emit(sink, session.Draft, session.Data{"revision": revision})
		emit(sink, session.Handoff, nil)
	}

	_, fenceSHA, err := service.Records.Fence(config)
	if err != nil {
		return Result{}, err
	}

	plan := service.plan(ticket, collected, request, guards, config, runtimeID, record, resources, components, totals, checks, state, fenceSHA, by)

	document := Document{
		Name:          config.Name,
		TicketName:    ticket.Name(),
		Run:           runtimeID.String(),
		ContextSHA256: collected.ContextSHA256,
		Record:        record.Record,
		By:            record.SetBy,
		Date:          record.SetAt.Format("2006-01-02"),
		Config:        config,
		Plan:          plan,
	}

	markdown, err := service.Records.Render(document)
	if err != nil {
		return Result{}, err
	}

	if err := service.Records.WriteRuntime(owner, runtimeID, markdown); err != nil {
		return Result{}, err
	}

	if err := service.Records.WriteRun(owner, runtimeID, service.run(ticket, guards, runtimeID, state)); err != nil {
		return Result{}, err
	}

	if err := service.Records.WriteLedger(owner, runtimeID.String(), service.ledger(ticket, guards, runtimeID, now)); err != nil {
		return Result{}, err
	}

	if state == Blocked {
		return Result{}, fmt.Errorf("%w: %s failed; the plan is written as blocked", ErrBlocked, strings.Join(failing, ", "))
	}

	return Result{Runtime: runtimeID.String(), Revision: revision, Path: runtimePath(runtimeID), State: state}, nil
}

func validate(request Request) error {
	if strings.TrimSpace(request.Ticket) == "" {
		return invalid("name the ticket to plan")
	}

	if len(request.Components) == 0 {
		return invalid("name at least one component with --component")
	}

	for _, component := range request.Components {
		if component.Path == "" || component.Image == "" || component.Port == 0 || component.ReadinessPath == "" || component.Build.Executable == "" {
			return invalid("component %q needs path, image, port, readiness and build", component.Path)
		}
	}

	if request.Limits.MaxDepth < 1 || request.Limits.MaxFiles < 1 || request.Limits.MaxBytes < 1 {
		return invalid("the collection limits must be positive")
	}

	if (request.CPUMillis < 0) || (request.MemoryMiB < 0) || (request.CPUMillis > 0) != (request.MemoryMiB > 0) {
		return invalid("override the envelope with both --cpu-millis and --memory-mib, or neither")
	}

	return nil
}

func invalid(format string, arguments ...any) error {
	return fmt.Errorf("%w: %s", inputs.ErrInvalid, fmt.Sprintf(format, arguments...))
}

// guardrails are the inputs a plan needs, as recorded.
type guardrails struct {
	envelope   harness.Envelope
	budget     finops.Budget
	choice     harness.Choice
	prices     *finops.Prices
	pricesAt   string
	forecasts  *finops.Forecasts
	workflow   *harness.Workflow
	workflowAt string
}

func (service Service) guardrails(owner project.Project) (guardrails, error) {
	var guards guardrails

	if err := load(service.Inputs, owner, inputs.KindEnvelope, &guards.envelope); err != nil {
		return guards, err
	}

	if err := load(service.Inputs, owner, inputs.KindBudget, &guards.budget); err != nil {
		return guards, err
	}

	if err := load(service.Inputs, owner, inputs.KindChoice, &guards.choice); err != nil {
		return guards, err
	}

	var prices finops.Prices
	mirror, err := loadOptional(service.Inputs, owner, inputs.KindPrices, &prices)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.prices, guards.pricesAt = &prices, mirror
	}

	var forecasts finops.Forecasts
	mirror, err = loadOptional(service.Inputs, owner, inputs.KindForecasts, &forecasts)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.forecasts = &forecasts
	}

	var workflow harness.Workflow
	mirror, err = loadOptional(service.Inputs, owner, inputs.KindWorkflow, &workflow)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.workflow, guards.workflowAt = &workflow, mirror
	}

	return guards, nil
}

// load reads a required input; a missing one is a missing prerequisite.
func load(store inputs.Store, owner project.Project, kind string, value any) error {
	record, err := store.Load(owner, kind)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(record.Value, value); err != nil {
		return fmt.Errorf("the %s input is unreadable: %w", kind, err)
	}

	return nil
}

// loadOptional reads an input nobody may have set, and says where its mirror
// is; an unset input leaves the value untouched and no mirror.
func loadOptional(store inputs.Store, owner project.Project, kind string, value any) (string, error) {
	record, err := store.Load(owner, kind)
	if errors.Is(err, application.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(record.Value, value); err != nil {
		return "", fmt.Errorf("the %s input is unreadable: %w", kind, err)
	}

	return orDefault(record.Mirror, "-"), nil
}

// ticketPathOf makes the ticket's path project-relative and checks it stays inside.
func ticketPathOf(owner project.Project, ticket string) (string, error) {
	relative := ticket
	if filepath.IsAbs(ticket) {
		inside, err := filepath.Rel(owner.Root, ticket)
		if err != nil || collection.Escapes(filepath.ToSlash(inside)) {
			return "", application.NotFound("the ticket %s is outside the project", ticket)
		}
		relative = inside
	}

	relative = filepath.ToSlash(filepath.Clean(relative))
	if collection.Escapes(relative) {
		return "", application.NotFound("the ticket %s is outside the project", ticket)
	}

	return relative, nil
}

func runtimeIDOf(override, slug string) (microruntime.Identifier, error) {
	text := override
	if text == "" {
		derived, err := collection.RuntimeID(slug)
		if err != nil {
			return microruntime.Identifier{}, invalid("%v; name one with --runtime", err)
		}
		text = derived
	}

	identifier, err := microruntime.ParseIdentifier(text)
	if err != nil {
		return microruntime.Identifier{}, invalid("runtime id %q: %v", text, err)
	}

	return identifier, nil
}

// revision is the next revision of the runtime, and its uid, kept across
// revisions. A fence that no longer validates still names its revision.
func (service Service) revision(owner project.Project, runtime microruntime.Identifier) (int, string, error) {
	existing, err := service.Records.ReadRuntime(owner, runtime)
	if errors.Is(err, application.ErrNotFound) {
		uid, err := service.IDs.ID()

		return 1, uid, err
	}
	if err != nil {
		return 0, "", err
	}

	config, err := connection.Parse(string(existing))
	if config.Revision < 1 {
		return 0, "", fmt.Errorf("the existing runtime.md of %s does not parse: %w", runtime, err)
	}

	return config.Revision + 1, config.RuntimeUID, nil
}

// bounds are the envelope the run takes: the owner's override, recorded as
// this runtime's input, or the recorded default.
func bounds(request Request, guards guardrails, runtime microruntime.Identifier, by string, now harness.Timestamp) (harness.Resources, run.InputValue) {
	if request.CPUMillis > 0 {
		resources := harness.Resources{CPUMillis: request.CPUMillis, MemoryMiB: request.MemoryMiB}

		return resources, run.InputValue{Record: "runtime-" + runtime.String() + "-input", SetBy: by, Via: inputs.Via, SetAt: now}
	}

	mirror := guards.envelope.Mirror

	return guards.envelope.Default, run.InputValue{Record: mirror.Record, SetBy: mirror.SetBy, Via: mirror.Via, SetAt: mirror.SetAt}
}

// resolved is one component, as named and as its manifests describe it.
type resolved struct {
	input        ComponentInput
	manifestName string
	technologies []run.Technology
	kind         string
	proposed     harness.Resources
}

// resolve inspects each component's manifests and sizes it, printing the
// inspection and the resulting row. Input wins over what a manifest says.
func (service Service) resolve(owner project.Project, components []ComponentInput, sink session.Sink) ([]resolved, error) {
	results := make([]resolved, 0, len(components))

	for _, input := range components {
		inspection, err := service.Inspector.Inspect(owner, input.Path)
		if err != nil {
			return nil, fmt.Errorf("component %s: %w", input.Path, err)
		}

		component := resolved{input: input, manifestName: orDefault(inspection.Manifest, "no manifest"), kind: orDefault(input.Kind, orDefault(inspection.Kind, manifest.Web))}

		technologies := inspection.Technologies
		if len(input.Technologies) > 0 {
			component.technologies = input.Technologies
			technologies = nil
			for _, technology := range input.Technologies {
				technologies = append(technologies, manifest.Technology{ID: technology.ID, Label: technology.Label})
			}
		} else {
			for _, technology := range technologies {
				component.technologies = append(component.technologies, run.Technology{ID: technology.ID, Label: technology.Label})
			}
		}

		component.proposed = manifest.Proposed(component.kind, technologies)

		emit(sink, session.Inspect, session.Data{"path": input.Path, "manifest": component.manifestName, "port": input.Port})
		emit(sink, session.Component, session.Data{"path": input.Path, "technology": manifest.Labels(technologies), "kind": component.kind, "cpuMillis": component.proposed.CPUMillis, "memoryMiB": component.proposed.MemoryMiB})

		results = append(results, component)
	}

	return results, nil
}

// crossChecks are the checks of the draft beyond the schema.
func crossChecks(components []resolved, bounded, ceiling harness.Resources) []run.Check {
	ports := make([]string, 0, len(components))
	distinct := map[int]bool{}
	for _, component := range components {
		ports = append(ports, strconv.Itoa(component.input.Port))
		distinct[component.input.Port] = true
	}
	portsCheck := run.Check{ID: "ports-distinct", Result: passed, Detail: text(strings.Join(ports, " · "))}
	if len(distinct) != len(components) {
		portsCheck.Result = failed
		portsCheck.Detail = text("ports repeat: " + strings.Join(ports, " · "))
	}

	images := make([]string, 0, len(components))
	imagesCheck := run.Check{ID: "images-resolved", Result: passed}
	for _, component := range components {
		images = append(images, component.input.Image)
		if !configuration.IsImageReference(component.input.Image) {
			imagesCheck.Result = failed
		}
	}
	imagesCheck.Detail = text(strings.Join(images, " · "))

	recipes := make([]string, 0, len(components))
	recipeCheck := run.Check{ID: "recipe-allowed", Result: passed}
	for _, component := range components {
		if err := component.input.Build.Validate(); err != nil {
			recipeCheck.Result = failed
			recipes = append(recipes, err.Error())

			continue
		}
		recipes = append(recipes, recipeName(component.input.Build))
	}
	recipeCheck.Detail = text(strings.Join(recipes, " · "))

	ceilingCheck := run.Check{ID: "within-ceiling", Result: passed, Detail: text(fmt.Sprintf("%sm of %sm · %s MiB of %s MiB", grouped(bounded.CPUMillis), grouped(ceiling.CPUMillis), grouped(bounded.MemoryMiB), grouped(ceiling.MemoryMiB)))}
	if bounded.CPUMillis > ceiling.CPUMillis || bounded.MemoryMiB > ceiling.MemoryMiB {
		ceilingCheck.Result = failed
	}

	return []run.Check{portsCheck, imagesCheck, recipeCheck, ceilingCheck}
}

func recipeName(build configuration.Build) string {
	if strings.Contains(build.Executable, "/") {
		return "committed script " + build.Executable
	}

	return strings.TrimSpace(build.Executable + " " + strings.Join(build.Args, " "))
}

func distinctPorts(components []resolved) int {
	ports := map[int]bool{}
	for _, component := range components {
		ports[component.input.Port] = true
	}

	return len(ports)
}

func collisionsOf(checks []run.Check) string {
	for _, check := range checks {
		if check.ID == "ports-distinct" && check.Result == failed {
			return "ports collide"
		}
	}

	return "no collisions"
}

func failingChecks(checks []run.Check) []string {
	var failing []string
	for _, check := range checks {
		if check.Result == failed {
			failing = append(failing, check.ID)
		}
	}

	return failing
}

// totalsOf sums the proposed shares against the envelope; headroom never
// goes below zero.
func totalsOf(components []resolved, bounded harness.Resources) run.Totals {
	var proposed harness.Resources
	for _, component := range components {
		proposed.CPUMillis += component.proposed.CPUMillis
		proposed.MemoryMiB += component.proposed.MemoryMiB
	}

	return run.Totals{
		Proposed: proposed,
		Bounded:  bounded,
		Headroom: harness.Resources{CPUMillis: max(0, bounded.CPUMillis-proposed.CPUMillis), MemoryMiB: max(0, bounded.MemoryMiB-proposed.MemoryMiB)},
	}
}

func (service Service) configuration(request Request, guards guardrails, owner project.Project, runtime microruntime.Identifier, uid string, revision int, resources harness.Resources, name string) configuration.Configuration {
	config := configuration.Configuration{
		Version:    1,
		RuntimeUID: uid,
		RuntimeID:  runtime.String(),
		ProjectID:  owner.ID.String(),
		Revision:   revision,
		Name:       name,
		CPUMillis:  resources.CPUMillis,
		MemoryMiB:  resources.MemoryMiB,
		Model:      guards.choice.Model,
		Effort:     guards.choice.Effort,
	}

	if guards.budget.Bounded.PerRunCents > 0 {
		budget := float64(guards.budget.Bounded.PerRunCents) / 100
		config.LLMBudgetUSD = &budget
	}

	for _, component := range request.Components {
		build := component.Build
		if build.Args == nil {
			build.Args = []string{}
		}
		config.Components = append(config.Components, configuration.Component{
			Path:          component.Path,
			Image:         component.Image,
			Port:          component.Port,
			ReadinessPath: component.ReadinessPath,
			Build:         &build,
		})
	}

	return config
}

func (service Service) plan(ticket collection.Ticket, collected collection.Result, request Request, guards guardrails, config configuration.Configuration, runtime microruntime.Identifier, record run.InputValue, resources harness.Resources, resolvedComponents []resolved, totals run.Totals, checks []run.Check, state, fenceSHA, by string) run.Plan {
	shares := run.SplitEnvelope(resources, len(resolvedComponents))

	components := make([]run.Component, 0, len(resolvedComponents))
	for index, component := range resolvedComponents {
		input := component.input
		inferredFrom := []string{}
		if component.manifestName != "no manifest" {
			inferredFrom = append(inferredFrom, path.Join(input.Path, component.manifestName))
		}

		components = append(components, run.Component{
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
			Executed:      run.Executed{CPUMillis: shares[index].CPUMillis, MemoryMiB: shares[index].MemoryMiB, Rule: run.EqualSplitRule},
		})
	}

	sources := make([]run.Source, 0, len(collected.Sources))
	for _, source := range collected.Sources {
		sources = append(sources, sourceRecord(source))
	}

	allocation := allocationOf(guards.budget, ticket)

	price := run.Price{RuntimePriced: false, Note: "the runtime has no price source; cpu and memory stay in units"}
	var modelName *string
	if guards.prices != nil {
		price.Record, price.Mirror, price.ReadOn = text(guards.prices.Mirror.Record), text(".langdrift/"+guards.pricesAt), text(guards.prices.ReadOn)
		for _, model := range guards.prices.Models {
			if model.ID == guards.choice.Model {
				modelName = text(model.Name)
			}
		}
	}

	envelope := guards.envelope

	return run.Plan{
		Version: 1,
		Run:     config.RuntimeID,
		State:   state,
		Ticket: run.Ticket{
			Path:    ticket.Path,
			Number:  optional(ticket.Number),
			Slug:    ticket.Slug,
			Title:   optional(ticket.Title),
			Status:  optional(ticket.Status),
			NamedBy: run.Actor{Person: by, Via: inputs.Via},
		},
		Collection: run.Collection{
			Agent:   nil,
			Limits:  run.Limits{MaxDepth: request.Limits.MaxDepth, MaxFiles: request.Limits.MaxFiles, MaxBytes: request.Limits.MaxBytes, Secrets: collection.RefusedWord, OutsideProject: collection.RefusedWord, Symlinks: collection.RefusedWord},
			Reached: run.Reached{Files: collected.Files, Depth: collected.Depth, Bytes: collected.Bytes},
			Context: run.Context{SHA256: collected.ContextSHA256, Rule: collection.ContextRule},
		},
		Sources: sources,
		Facts:   []run.Fact{},
		Inputs: run.Inputs{
			CPUMillis:        run.InputValue{Value: resources.CPUMillis, Record: record.Record, SetBy: record.SetBy, Via: record.Via, SetAt: record.SetAt},
			MemoryMiB:        run.InputValue{Value: resources.MemoryMiB, Record: record.Record, SetBy: record.SetBy, Via: record.Via, SetAt: record.SetAt},
			ModelBudgetCents: run.InputValue{Value: allocation, Record: guards.budget.Mirror.Record, SetBy: guards.budget.Mirror.SetBy, Via: guards.budget.Mirror.Via, SetAt: guards.budget.Mirror.SetAt},
			Ceiling:          run.Ceiling{CPUMillis: envelope.PerExecution.CPUMillis, MemoryMiB: envelope.PerExecution.MemoryMiB, Record: envelope.Mirror.Record, SetBy: envelope.Mirror.SetBy, Via: envelope.Mirror.Via, SetAt: envelope.Mirror.SetAt},
			NotRuntimeFields: nonNil(envelope.NotRuntimeFields),
		},
		Discrepancies: []run.Discrepancy{},
		Retrieval: run.Retrieval{
			Method:        "vector+structured",
			Criteria:      []string{"same-component-path", "same-stack", "measured-resources", "measured-tokens"},
			KeepAtMatches: 2,
			Considered:    []run.Considered{},
			Note:          "retrieval not run yet",
		},
		Sizing: run.Sizing{Basis: "input", MeasuredRunsCited: []string{}, Note: "no episode retrieved yet · the proposed shares follow each component's kind"},
		Suggestion: run.Suggestion{
			By:           "input",
			Model:        guards.choice.Model,
			ModelName:    modelName,
			Effort:       guards.choice.Effort,
			Reasons:      []string{},
			AuthorizedBy: "pending",
		},
		Components:      components,
		Totals:          totals,
		Materialization: run.Materialization{Deployments: len(components), ReplicasEach: 1},
		Price:           price,
		Checks:          checks,
		Runtime:         run.RuntimeRef{Path: runtimePath(runtime), Revision: config.Revision, SHA256: fenceSHA, State: state},
		ContractGaps:    append([]string{}, run.ContractGaps...),
	}
}

func allocationOf(budget finops.Budget, ticket collection.Ticket) int {
	for _, entry := range budget.Allocations {
		if entry.Ticket == ticket.Reference() {
			return entry.Cents
		}
	}

	return 0
}

// run is the execution record of a planned runtime: nothing has run yet.
func (service Service) run(ticket collection.Ticket, guards guardrails, runtime microruntime.Identifier, state string) run.Run {
	record := run.Run{
		Version:  1,
		ID:       runtime.String(),
		Kind:     "ticket",
		Ticket:   text(ticket.Path),
		Runtime:  runtime.String(),
		State:    "planned",
		Files:    run.Files{Plan: text("runtime.md"), Memory: "run.md", Ledger: ".langdrift/finops/ledger/" + runtime.String() + ".ledger.json", Evidence: "evidence/"},
		Laps:     []run.LapSummary{},
		Stages:   []run.StageState{},
		Delivery: run.Delivery{Checks: []run.CheckResult{}},
	}

	if state == Blocked {
		record.State = Blocked
	}

	record.HumanInTheLoop = run.HumanInTheLoop{Chosen: []string{}}
	if guards.workflow == nil {
		return record
	}

	workflow := guards.workflow
	record.HumanInTheLoop = run.HumanInTheLoop{Record: text(workflow.Mirror.Record), Mirror: text(".langdrift/" + guards.workflowAt), Chosen: nonNil(workflow.Defaults.HumanAt)}

	for _, stage := range workflow.Stages {
		record.Stages = append(record.Stages, run.StageState{ID: stage.ID, State: "pending", Human: isHuman(stage, workflow.Defaults)})
	}

	return record
}

func isHuman(stage harness.Stage, defaults harness.Defaults) bool {
	if stage.Human == "required" || stage.Human == "owner-only" {
		return true
	}

	for _, id := range defaults.HumanAt {
		if id == stage.ID {
			return true
		}
	}

	return false
}

// ledger opens the runtime's ledger with its planning line: nothing paid,
// since no collection agent ran.
func (service Service) ledger(ticket collection.Ticket, guards guardrails, runtime microruntime.Identifier, now harness.Timestamp) finops.Ledger {
	ledger := finops.Ledger{
		Mirror:   harness.Mirror{Of: inputs.Via, Record: "ledger-" + runtime.String(), SetBy: "ex", Via: inputs.Via, SetAt: now, ReadOnly: true, Note: inputs.MirrorNote},
		Version:  1,
		Run:      runtime.String(),
		Ticket:   text(ticket.Reference()),
		State:    "provisional",
		Currency: guards.budget.Currency,
		Bounded:  finops.LedgerBounds{PerRunCents: guards.budget.Bounded.PerRunCents, TicketAllocationCents: allocationOf(guards.budget, ticket), PeriodLeftCents: guards.budget.Bounded.PeriodCents},
		Planning: finops.PlanningLine{Model: nil, Note: "no collection agent ran; the plan was built by the deterministic shell, nothing paid", Cents: finops.PlanningCents{Billed: number(0)}},
		Laps:     []finops.LapLine{},
		Totals:   finops.LedgerTotals{Cents: finops.TotalCents{Billed: number(0)}, PeriodLeftAfterCents: guards.budget.Bounded.PeriodCents},
	}

	if guards.prices != nil {
		ledger.PriceSource = text(".langdrift/" + guards.pricesAt)
	}
	if guards.forecasts != nil {
		ledger.ForecastSource = text(".langdrift/" + inputs.ForecastsMirror)
	}

	return ledger
}

func sourceRecord(source collection.Source) run.Source {
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

	return record
}

func isGuardrail(source collection.Source) bool {
	return source.State == collection.Read && (strings.HasSuffix(source.Path, ".gr.md") || strings.Contains(source.Path, "guardrail"))
}

// criteria are the criterion ids a guardrail names, in order, once each.
func criteria(content []byte) []string {
	var found []string
	seen := map[string]bool{}
	for _, id := range criterionID.FindAllString(string(content), -1) {
		if seen[id] {
			continue
		}
		seen[id] = true
		found = append(found, id)
	}

	return found
}

func linkNames(collected collection.Result) []string {
	var names []string
	for _, source := range collected.Sources {
		if source.Depth != 1 {
			continue
		}
		names = append(names, path.Base(source.Path))
	}

	return names
}

func judgeReason(request Request) string {
	if request.NoJudge {
		return "--no-judge; verdicts, facts and the suggestion are the inputs' own"
	}

	return "no judge is configured yet; verdicts, facts and the suggestion are the inputs' own"
}

func componentName(componentPath string) string {
	if componentPath == "." {
		return "root"
	}

	return path.Base(componentPath)
}

func runtimePath(runtime microruntime.Identifier) string {
	return microruntime.ConfigurationDirectory(runtime).String() + "/runtime.md"
}

// grouped prints a figure with thousands separators, as the checks do.
func grouped(value int) string {
	digits := strconv.Itoa(value)
	if len(digits) <= 3 {
		return digits
	}

	var parts []string
	for len(digits) > 3 {
		parts = append([]string{digits[len(digits)-3:]}, parts...)
		digits = digits[:len(digits)-3]
	}
	parts = append([]string{digits}, parts...)

	return strings.Join(parts, ",")
}

func emit(sink session.Sink, id string, data session.Data) {
	sink.Emit(session.Event{Session: session.PlanSession, ID: id, Data: data})
}

func text(value string) *string { return &value }

func optional(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func number(value int) *int { return &value }

func size(value int64) *int64 { return &value }

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func nonNil[Item any](items []Item) []Item {
	if items == nil {
		return []Item{}
	}

	return items
}
