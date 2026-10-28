package finops

import "github.com/jonatassales/tao/packages/ex/src/domain/harness"

// Prices is a dated price list mirror, in cents per million tokens.
type Prices struct {
	Mirror       harness.Mirror `json:"mirror"`
	Version      int            `json:"version"`
	Illustrative bool           `json:"illustrative,omitempty"`
	Note         string         `json:"note"`
	ReadOn       string         `json:"readOn"`
	Currency     string         `json:"currency"`
	Unit         string         `json:"unit"`
	Cache        string         `json:"cache"`
	Models       []ModelPrice   `json:"models"`
}

// ModelPrice is one model's list price for input and output tokens.
type ModelPrice struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	InputCents  int    `json:"inputCents"`
	OutputCents int    `json:"outputCents"`
}
