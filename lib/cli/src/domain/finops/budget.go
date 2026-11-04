// Package finops holds the money records: the budget, the gates, the prices
// and the forecasts the CLI mirrors, and the ledger of each runtime. Money
// formats belong to FinOps, not to the Microruntimes contracts.
package finops

import "github.com/jonatassales/tao/lib/cli/src/domain/harness"

// Budget is the team budget mirror for one period.
type Budget struct {
	Mirror           harness.Mirror `json:"mirror"`
	Version          int            `json:"version"`
	Status           string         `json:"status"`
	Owner            string         `json:"owner"`
	Currency         string         `json:"currency"`
	Period           string         `json:"period"`
	Bounded          BudgetBounds   `json:"bounded"`
	Allocations      []Allocation   `json:"allocations"`
	UnallocatedCents int            `json:"unallocatedCents"`
}

// BudgetBounds are the period's amount, the per-run cap and the margin a
// reservation adds to an estimate.
type BudgetBounds struct {
	PeriodCents          int `json:"periodCents"`
	PerRunCents          int `json:"perRunCents"`
	ReserveMarginPercent int `json:"reserveMarginPercent"`
}

// Allocation is the share of the period one ticket may spend.
type Allocation struct {
	Ticket string `json:"ticket"`
	Cents  int    `json:"cents"`
}
