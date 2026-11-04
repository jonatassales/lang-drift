package finops

import "github.com/jonatassales/tao/lib/cli/src/domain/harness"

// Forecasts is the mirror of the tokens a run is expected to use by effort,
// used when no earlier run measured tokens for the ticket.
type Forecasts struct {
	Mirror       harness.Mirror `json:"mirror"`
	Version      int            `json:"version"`
	Status       string         `json:"status"`
	Illustrative bool           `json:"illustrative,omitempty"`
	Note         string         `json:"note"`
	Unit         string         `json:"unit"`
	ByEffort     ByEffort       `json:"byEffort"`
	// Planning is what the planner's own calls are expected to use; absent
	// until the owner records it.
	Planning *TokensThousands `json:"planning,omitempty"`
}

// ByEffort is the forecast for each effort. A struct, not a map, so the
// efforts keep their order in the mirror.
type ByEffort struct {
	Low    TokensThousands `json:"low"`
	Medium TokensThousands `json:"medium"`
	High   TokensThousands `json:"high"`
}

// TokensThousands counts input and output tokens in thousands.
type TokensThousands struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}
