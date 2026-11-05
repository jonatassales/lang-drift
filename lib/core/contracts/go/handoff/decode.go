package handoff

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Decode reads one hand-off record, refusing a field its shape does not
// define and a text that holds more than one record.
func Decode[Record any](text string) (Record, error) {
	var record Record

	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&record); err != nil {
		return record, fmt.Errorf("handoff: %w", err)
	}

	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return record, errors.New("handoff: the text must hold one record")
	}

	return record, nil
}
