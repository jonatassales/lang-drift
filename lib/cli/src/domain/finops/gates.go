package finops

import "github.com/jonatassales/tao/lib/cli/src/domain/harness"

// Gates is the mirror of the money gates a run passes.
type Gates struct {
	Mirror  harness.Mirror `json:"mirror"`
	Version int            `json:"version"`
	Status  string         `json:"status"`
	Gates   []Gate         `json:"gates"`
}

// Gate is one deterministic control: when it runs, what blocks, what it does
// and how it holds today.
type Gate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Criterion   string `json:"criterion"`
	When        string `json:"when"`
	BlocksIf    string `json:"blocksIf"`
	Action      string `json:"action"`
	Enforcement string `json:"enforcement"`
}

// DefaultGates are the six gates of the reference policy, FG-01 to FG-06.
func DefaultGates() []Gate {
	return []Gate{
		{ID: "FG-01", Name: "price-source-required", Criterion: "FIN-02", When: "before-authorization", BlocksIf: "the plan has paid work and no dated price source covers it", Action: "block", Enforcement: "planned"},
		{ID: "FG-02", Name: "budget-required", Criterion: "FIN-03", When: "before-authorization", BlocksIf: "the ticket has no allocation recorded on the server", Action: "block", Enforcement: "planned"},
		{ID: "FG-03", Name: "reservation-within-bounds", Criterion: "FIN-03", When: "before-authorization", BlocksIf: "estimate plus margin exceeds the per-run cap, the ticket allocation or what is left of the period", Action: "pause", Enforcement: "planned"},
		{ID: "FG-04", Name: "reservation-exhausted", Criterion: "FIN-03", When: "during-run", BlocksIf: "observed spend reaches the reservation", Action: "pause-paid-calls", Enforcement: "planned"},
		{ID: "FG-05", Name: "overspend-acknowledged", Criterion: "FIN-03", When: "before-authorization", BlocksIf: "an earlier run of the ticket has overspend the owner has not acknowledged", Action: "require-owner", Enforcement: "planned"},
		{ID: "FG-06", Name: "billing-pending", Criterion: "FIN-02", When: "after-run", BlocksIf: "the provider has not billed the run yet", Action: "keep-provisional", Enforcement: "guidance"},
	}
}
