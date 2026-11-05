package stages

import (
	"encoding/json"
	"errors"
)

// typed names the field whose value has the wrong type, when decoding the
// report failed on one.
func typed(err error) error {
	if err == nil {
		return nil
	}

	var mistyped *json.UnmarshalTypeError
	if !errors.As(err, &mistyped) {
		return FieldError{
			Field: "report",
			Rule:  err.Error(),
		}
	}

	return FieldError{
		Field: mistyped.Field,
		Rule:  "must be a " + mistyped.Type.String(),
	}
}
