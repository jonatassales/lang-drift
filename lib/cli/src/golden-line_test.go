package main

import (
	"encoding/json"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// goldenLine is one line of a frozen golden session: its line id, its kind
// and its exact text.
type goldenLine struct {
	id   string
	kind session.Kind
	text string
}

// UnmarshalJSON decodes one line from a frozen golden JSON array.
func (line *goldenLine) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID   string       `json:"id"`
		Kind session.Kind `json:"kind"`
		Text string       `json:"text"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	line.id = raw.ID
	line.kind = raw.Kind
	line.text = raw.Text

	return nil
}
