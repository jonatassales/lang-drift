package stages

import (
	"encoding/json"
	"slices"
)

// null is how JSON writes a value that is not there.
const null = "null"

// readReport reads one JSON object as a report holding exactly the named
// fields, each present and none null, into the value. Unknown is never zero,
// so a null count or flag is refused, not read as 0 or false.
func readReport(text string, fields []string, into any) error {
	var object map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &object); err != nil {
		return FieldError{
			Field: "report",
			Rule:  "must be a JSON object: " + err.Error(),
		}
	}

	for _, field := range fields {
		value, present := object[field]
		if !present {
			return FieldError{
				Field: field,
				Rule:  "is missing",
			}
		}

		if string(value) == null {
			return FieldError{
				Field: field,
				Rule:  "must not be null",
			}
		}
	}

	for field := range object {
		if !slices.Contains(fields, field) {
			return FieldError{
				Field: field,
				Rule:  "is not a field of the report",
			}
		}
	}

	return typed(json.Unmarshal([]byte(text), into))
}
