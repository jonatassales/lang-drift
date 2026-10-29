package plan

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// Guardrails are the inputs a plan and a run need, as recorded: the
// required ones, and the optional ones with where their mirror is.
type Guardrails struct {
	Envelope   harness.Envelope
	Budget     finops.Budget
	Choice     harness.Choice
	Prices     *finops.Prices
	PricesAt   string
	Forecasts  *finops.Forecasts
	Workflow   *harness.Workflow
	WorkflowAt string
}

// LoadGuardrails reads the inputs; a missing required one is a missing
// prerequisite that names itself.
func LoadGuardrails(store inputs.Store, owner project.Project) (Guardrails, error) {
	var guards Guardrails

	if err := load(store, owner, inputs.KindEnvelope, &guards.Envelope); err != nil {
		return guards, err
	}

	if err := load(store, owner, inputs.KindBudget, &guards.Budget); err != nil {
		return guards, err
	}

	if err := load(store, owner, inputs.KindChoice, &guards.Choice); err != nil {
		return guards, err
	}

	var prices finops.Prices
	mirror, err := loadOptional(store, owner, inputs.KindPrices, &prices)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.Prices, guards.PricesAt = &prices, mirror
	}

	var forecasts finops.Forecasts
	mirror, err = loadOptional(store, owner, inputs.KindForecasts, &forecasts)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.Forecasts = &forecasts
	}

	var workflow harness.Workflow
	mirror, err = loadOptional(store, owner, inputs.KindWorkflow, &workflow)
	if err != nil {
		return guards, err
	}
	if mirror != "" {
		guards.Workflow, guards.WorkflowAt = &workflow, mirror
	}

	return guards, nil
}

// Price is a model's list price, when the current price list covers it.
func (guards Guardrails) Price(model string) (finops.ModelPrice, bool) {
	if guards.Prices == nil {
		return finops.ModelPrice{}, false
	}

	for _, entry := range guards.Prices.Models {
		if entry.ID == model {
			return entry, true
		}
	}

	return finops.ModelPrice{}, false
}

// Allocation is the ticket's share of the period, or zero.
func (guards Guardrails) Allocation(ticket collection.Ticket) int {
	return guards.AllocationOf(ticket.Reference())
}

// AllocationOf is a ticket reference's share of the period, or zero.
func (guards Guardrails) AllocationOf(reference string) int {
	for _, entry := range guards.Budget.Allocations {
		if entry.Ticket == reference {
			return entry.Cents
		}
	}

	return 0
}

// Planner is the model the planner's own calls use.
func (guards Guardrails) Planner() string {
	return orDefault(guards.Choice.Planner, guards.Choice.Model)
}

// Forecast is the tokens an effort is expected to use, when forecasts exist.
func (guards Guardrails) Forecast(effort string) (finops.TokensThousands, bool) {
	if guards.Forecasts == nil {
		return finops.TokensThousands{}, false
	}

	switch effort {
	case "low":
		return guards.Forecasts.ByEffort.Low, true
	case "medium":
		return guards.Forecasts.ByEffort.Medium, true
	case "high":
		return guards.Forecasts.ByEffort.High, true
	}

	return finops.TokensThousands{}, false
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
