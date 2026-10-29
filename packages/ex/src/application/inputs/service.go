package inputs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// Service records the inputs. Each operation validates its request, keeps the
// record privately with the one it replaces, and mirrors it where a mirror
// exists.
type Service struct {
	Projects application.ProjectOpener
	Store    Store
	Mirrors  Mirrors

	// Injectable for tests; nil uses the real sources.
	IDs   application.IDSource
	Clock application.Clock
}

// DefaultBy is who sets an input, and who a run acts for, unless named.
const DefaultBy = "owner"

// Defaults every record takes when the request leaves them empty.
const (
	defaultOwner  = "team"
	defaultStatus = "draft"
	recordVersion = 1
)

var (
	periodPattern = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)
	datePattern   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	slugPattern   = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
)

// EnvelopeRequest sets the hardware envelope.
type EnvelopeRequest struct {
	Project          application.ProjectSelection
	By               string
	Owner            string
	Status           string
	PerExecution     harness.Resources
	Default          harness.Resources
	NotRuntimeFields []string
}

// SetEnvelope records the hardware envelope and mirrors it.
func (service Service) SetEnvelope(_ context.Context, request EnvelopeRequest) (Result, error) {
	if err := validateEnvelope(request); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindEnvelope, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}

	envelope := harness.Envelope{
		Mirror:           stamp.mirror(),
		Version:          recordVersion,
		Status:           orDefault(request.Status, defaultStatus),
		Owner:            orDefault(request.Owner, defaultOwner),
		PerExecution:     request.PerExecution,
		Default:          request.Default,
		NotRuntimeFields: nonNil(request.NotRuntimeFields),
	}

	return service.record(owner, stamp, envelope, EnvelopeMirror)
}

func validateEnvelope(request EnvelopeRequest) error {
	if !positiveResources(request.PerExecution) || !positiveResources(request.Default) {
		return invalid("the envelope's cpuMillis and memoryMiB must be positive")
	}

	if request.Default.CPUMillis > request.PerExecution.CPUMillis || request.Default.MemoryMiB > request.PerExecution.MemoryMiB {
		return invalid("the default envelope must stay inside the per-execution ceiling")
	}

	return nil
}

// WorkflowRequest sets the delivery stages. Mirror and Version are the
// service's to set; the rest is the owner's.
type WorkflowRequest struct {
	Project  application.ProjectSelection
	By       string
	Workflow harness.Workflow
}

// SetWorkflow records the delivery stages and mirrors them.
func (service Service) SetWorkflow(_ context.Context, request WorkflowRequest) (Result, error) {
	workflow := request.Workflow
	workflow.Status = orDefault(workflow.Status, defaultStatus)
	workflow.Owner = orDefault(workflow.Owner, defaultOwner)
	if len(workflow.Efforts) == 0 {
		workflow.Efforts = append([]string{}, harness.Efforts...)
	}
	workflow.Defaults.HumanAt = nonNil(workflow.Defaults.HumanAt)
	for index := range workflow.Stages {
		workflow.Stages[index].Skills = nonNil(workflow.Stages[index].Skills)
	}

	if err := validateWorkflow(workflow); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindWorkflow, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}

	workflow.Mirror = stamp.mirror()
	workflow.Version = recordVersion

	return service.record(owner, stamp, workflow, WorkflowMirror)
}

func validateWorkflow(workflow harness.Workflow) error {
	if len(workflow.Stages) == 0 {
		return invalid("a workflow needs at least one stage")
	}

	ids := map[string]bool{}
	for _, stage := range workflow.Stages {
		if stage.ID == "" || ids[stage.ID] {
			return invalid("stage ids must be present and distinct; %q is not", stage.ID)
		}
		ids[stage.ID] = true

		if !validHuman(stage.Human) {
			return invalid("stage %s: human must be required, optional or owner-only, not %q", stage.ID, stage.Human)
		}
	}

	for _, id := range workflow.Defaults.HumanAt {
		if !ids[id] {
			return invalid("defaults.humanAt names %q, which is not a stage", id)
		}
	}

	for _, effort := range workflow.Efforts {
		if !harness.ValidEffort(effort) {
			return invalid("effort %q is not low, medium or high", effort)
		}
	}

	return nil
}

func validHuman(text string) bool {
	return text == harness.HumanRequired || text == harness.HumanOptional || text == harness.OwnerOnly
}

// ModelRequest sets the default model and effort, and the planner's model.
type ModelRequest struct {
	Project application.ProjectSelection
	By      string
	Model   string
	Effort  string
	Planner string
}

// SetModel records the default choice. It has no mirror: it blocks nothing by itself.
func (service Service) SetModel(_ context.Context, request ModelRequest) (Result, error) {
	if !plain(request.Model) {
		return Result{}, invalid("a model id is required")
	}

	if !harness.ValidEffort(request.Effort) {
		return Result{}, invalid("effort must be low, medium or high, not %q", request.Effort)
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindChoice, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}

	choice := harness.Choice{Model: request.Model, Effort: request.Effort, Planner: orDefault(request.Planner, request.Model)}

	return service.record(owner, stamp, choice, "")
}

// DriverRequest sets the coding agent the stages run through.
type DriverRequest struct {
	Project application.ProjectSelection
	By      string
	Driver  harness.Driver
}

// Defaults of a driver the owner left open.
const (
	defaultDriverTimeoutSeconds = 1800
	defaultUnmeasured           = "pause"
)

// SetDriver records the driver. It has no mirror: it blocks nothing by itself.
func (service Service) SetDriver(_ context.Context, request DriverRequest) (Result, error) {
	driver := request.Driver
	driver.Args = nonNil(driver.Args)
	driver.PassEnv = nonNil(driver.PassEnv)
	driver.Unmeasured = orDefault(driver.Unmeasured, defaultUnmeasured)
	if driver.TimeoutSeconds == 0 {
		driver.TimeoutSeconds = defaultDriverTimeoutSeconds
	}

	if err := validateDriver(driver); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindDriver, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}

	return service.record(owner, stamp, driver, "")
}

func validateDriver(driver harness.Driver) error {
	if !plain(driver.Executable) {
		return invalid("the driver needs an executable")
	}

	for _, argument := range driver.Args {
		if strings.ContainsAny(argument, "\n\r\x00") {
			return invalid("a driver argument holds a control character")
		}
	}

	if driver.TimeoutSeconds < 1 {
		return invalid("the driver's timeout must be positive")
	}

	if driver.Unmeasured != "pause" && driver.Unmeasured != "allow" {
		return invalid("unmeasured must be pause or allow, not %q", driver.Unmeasured)
	}

	return nil
}

// BudgetRequest sets the team budget for one period.
type BudgetRequest struct {
	Project     application.ProjectSelection
	By          string
	Owner       string
	Status      string
	Currency    string
	Period      string
	Bounded     finops.BudgetBounds
	Allocations []finops.Allocation
}

// SetBudget records the budget, computes what stays unallocated, and mirrors it.
func (service Service) SetBudget(_ context.Context, request BudgetRequest) (Result, error) {
	allocated, err := validateBudget(request)
	if err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindBudget+"-"+request.Period, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}
	stamp.kind = KindBudget

	budget := finops.Budget{
		Mirror:           stamp.mirror(),
		Version:          recordVersion,
		Status:           orDefault(request.Status, defaultStatus),
		Owner:            orDefault(request.Owner, defaultOwner),
		Currency:         orDefault(request.Currency, "USD"),
		Period:           request.Period,
		Bounded:          request.Bounded,
		Allocations:      nonNil(request.Allocations),
		UnallocatedCents: request.Bounded.PeriodCents - allocated,
	}

	return service.record(owner, stamp, budget, BudgetMirror)
}

func validateBudget(request BudgetRequest) (int, error) {
	if !periodPattern.MatchString(request.Period) {
		return 0, invalid("the period must be a month, YYYY-MM, not %q", request.Period)
	}

	bounds := request.Bounded
	if bounds.PeriodCents < 0 || bounds.PerRunCents < 0 || bounds.PerRunCents > bounds.PeriodCents {
		return 0, invalid("the per-run cap must stay between zero and the period's cents")
	}

	if bounds.ReserveMarginPercent < 0 || bounds.ReserveMarginPercent > 100 {
		return 0, invalid("the reserve margin must be a percentage between 0 and 100")
	}

	allocated := 0
	tickets := map[string]bool{}
	for _, allocation := range request.Allocations {
		if allocation.Ticket == "" || tickets[allocation.Ticket] {
			return 0, invalid("allocations name each ticket once; %q does not", allocation.Ticket)
		}
		tickets[allocation.Ticket] = true

		if allocation.Cents < 0 {
			return 0, invalid("the allocation of %s is negative", allocation.Ticket)
		}

		allocated += allocation.Cents
	}

	if allocated > bounds.PeriodCents {
		return 0, invalid("the allocations exceed the period's %d cents", bounds.PeriodCents)
	}

	return allocated, nil
}

// GatesRequest sets the money gates.
type GatesRequest struct {
	Project application.ProjectSelection
	By      string
	Status  string
	Gates   []finops.Gate
}

// SetGates records the gates and mirrors them.
func (service Service) SetGates(_ context.Context, request GatesRequest) (Result, error) {
	if err := validateGates(request.Gates); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindGates, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}

	gates := finops.Gates{
		Mirror:  stamp.mirror(),
		Version: recordVersion,
		Status:  orDefault(request.Status, defaultStatus),
		Gates:   request.Gates,
	}

	return service.record(owner, stamp, gates, GatesMirror)
}

func validateGates(gates []finops.Gate) error {
	if len(gates) == 0 {
		return invalid("at least one gate is required")
	}

	ids := map[string]bool{}
	for _, gate := range gates {
		if gate.ID == "" || ids[gate.ID] {
			return invalid("gate ids must be present and distinct; %q is not", gate.ID)
		}
		ids[gate.ID] = true

		if gate.Action == "" {
			return invalid("gate %s has no action", gate.ID)
		}
	}

	return nil
}

// PricesRequest sets a dated price list.
type PricesRequest struct {
	Project application.ProjectSelection
	By      string
	Slug    string
	Prices  finops.Prices
}

// PricesUnit is the only unit a price list may use.
const PricesUnit = "cents-per-million-tokens"

// SetPrices records a price list as the current one and mirrors it under the
// next free number.
func (service Service) SetPrices(_ context.Context, request PricesRequest) (Result, error) {
	prices := request.Prices
	prices.Currency = orDefault(prices.Currency, "USD")
	prices.Cache = orDefault(prices.Cache, "uncached")

	if err := validatePrices(request.Slug, prices); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	number, err := service.Mirrors.NextPriceNumber(owner)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp(KindPrices, orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}
	stamp.id = KindPrices + "-" + number

	prices.Mirror = stamp.mirror()
	prices.Version = recordVersion

	return service.record(owner, stamp, prices, PricesDirectory+"/"+number+"-"+request.Slug+".prices.json")
}

func validatePrices(slug string, prices finops.Prices) error {
	if !slugPattern.MatchString(slug) {
		return invalid("the price list's slug must be lowercase letters, digits and dashes, not %q", slug)
	}

	if prices.Unit != PricesUnit {
		return invalid("prices are in %s, not %q", PricesUnit, prices.Unit)
	}

	if !datePattern.MatchString(prices.ReadOn) {
		return invalid("readOn must be a date, YYYY-MM-DD, not %q", prices.ReadOn)
	}

	if len(prices.Models) == 0 {
		return invalid("a price list needs at least one model")
	}

	ids := map[string]bool{}
	for _, model := range prices.Models {
		if !plain(model.ID) || ids[model.ID] {
			return invalid("model ids must be present and distinct; %q is not", model.ID)
		}
		ids[model.ID] = true

		if model.InputCents < 0 || model.OutputCents < 0 {
			return invalid("model %s has a negative price", model.ID)
		}
	}

	return nil
}

// ForecastsRequest sets the token forecasts.
type ForecastsRequest struct {
	Project   application.ProjectSelection
	By        string
	Forecasts finops.Forecasts
}

// ForecastsUnit is the only unit a forecast may use.
const ForecastsUnit = "thousands-of-tokens"

// SetForecasts records the forecasts and mirrors them.
func (service Service) SetForecasts(_ context.Context, request ForecastsRequest) (Result, error) {
	forecasts := request.Forecasts
	forecasts.Status = orDefault(forecasts.Status, defaultStatus)

	if err := validateForecasts(forecasts); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	stamp, err := service.stamp("forecast", orDefault(request.By, DefaultBy))
	if err != nil {
		return Result{}, err
	}
	stamp.kind = KindForecasts

	forecasts.Mirror = stamp.mirror()
	forecasts.Version = recordVersion

	return service.record(owner, stamp, forecasts, ForecastsMirror)
}

func validateForecasts(forecasts finops.Forecasts) error {
	if forecasts.Unit != ForecastsUnit {
		return invalid("forecasts are in %s, not %q", ForecastsUnit, forecasts.Unit)
	}

	for _, tokens := range []finops.TokensThousands{forecasts.ByEffort.Low, forecasts.ByEffort.Medium, forecasts.ByEffort.High} {
		if tokens.Input <= 0 || tokens.Output <= 0 {
			return invalid("every effort needs positive input and output tokens")
		}
	}

	if forecasts.Planning != nil && (forecasts.Planning.Input <= 0 || forecasts.Planning.Output <= 0) {
		return invalid("the planning forecast needs positive input and output tokens")
	}

	return nil
}

// stamp is the provenance of one record: its kind and id, who set it and when.
type stamp struct {
	kind  string
	id    string
	by    string
	setAt harness.Timestamp
}

// stamp draws a record id from the prefix and a fresh identifier, and the
// moment from the clock.
func (service Service) stamp(prefix, by string) (stamp, error) {
	raw, err := service.IDs.ID()
	if err != nil {
		return stamp{}, err
	}

	if len(raw) < 4 {
		return stamp{}, fmt.Errorf("identifier %q is too short for a record id", raw)
	}

	return stamp{kind: prefix, id: prefix + "-" + raw[:4], by: by, setAt: harness.At(service.Clock.Now())}, nil
}

func (stamp stamp) mirror() harness.Mirror {
	return harness.Mirror{Of: Via, Record: stamp.id, SetBy: stamp.by, Via: Via, SetAt: stamp.setAt, ReadOnly: true, Note: MirrorNote}
}

// record stores the value with the record it replaces and writes its mirror.
func (service Service) record(owner project.Project, stamp stamp, value any, mirror string) (Result, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return Result{}, err
	}

	record := Record{Kind: stamp.kind, ID: stamp.id, SetBy: stamp.by, Via: Via, SetAt: stamp.setAt, Value: encoded, Mirror: mirror}

	previous, err := service.Store.Load(owner, stamp.kind)
	switch {
	case errors.Is(err, application.ErrNotFound):
		// The first record of its kind.

	case err != nil:
		return Result{}, err

	default:
		record.Previous = &previous
	}

	if err := service.Store.Save(owner, record); err != nil {
		return Result{}, err
	}

	if mirror != "" {
		if err := service.Mirrors.WriteMirror(owner, mirror, value); err != nil {
			return Result{}, err
		}
	}

	return Result{Kind: stamp.kind, Record: stamp.id, Mirror: mirror, Status: "set"}, nil
}

func positiveResources(resources harness.Resources) bool {
	return resources.CPUMillis > 0 && resources.MemoryMiB > 0
}

func plain(text string) bool {
	return strings.TrimSpace(text) != "" && !strings.ContainsAny(text, " \t\r\n\x00")
}

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
