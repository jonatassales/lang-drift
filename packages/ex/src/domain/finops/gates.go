package finops

import "github.com/jonatassales/tao/packages/ex/src/domain/harness"

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
