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

// price is a model's list price, when the current price list covers it.
func (guards guardrails) price(model string) (finops.ModelPrice, bool) {
	if guards.prices == nil {
		return finops.ModelPrice{}, false
	}

	for _, entry := range guards.prices.Models {
		if entry.ID == model {
			return entry, true
		}
	}

	return finops.ModelPrice{}, false
}

// allocation is the ticket's share of the period, or zero.
func (guards guardrails) allocation(ticket collection.Ticket) int {
	for _, entry := range guards.budget.Allocations {
		if entry.Ticket == ticket.Reference() {
			return entry.Cents
		}
	}

	return 0
}

// planner is the model the planner's own calls use.
func (guards guardrails) planner() string {
	return orDefault(guards.choice.Planner, guards.choice.Model)
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
