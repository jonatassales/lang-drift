package judgment

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// decodeStrictly reads one JSON object into the value, refusing any field
// the answer does not define and anything after the object.
func decodeStrictly(text string, into any) error {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(into); err != nil {
		return FieldError{
			Field: "answer",
			Rule:  "must be the kind's JSON object: " + err.Error(),
		}
	}

	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return FieldError{
			Field: "answer",
			Rule:  "must hold one JSON object",
		}
	}

	return nil
}
