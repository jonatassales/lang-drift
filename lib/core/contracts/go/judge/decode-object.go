package judge

import (
	"encoding/json"
	"errors"
	"slices"
)

// decodeObject reads one JSON object into the value, as the schemas read it:
// every field listed is present, no other field is, and each value has the
// type the struct gives it. The path names the object in a refusal; the
// answer itself has an empty path.
func decodeObject(raw json.RawMessage, path string, fields []string, into any) error {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return refuse(objectName(path), "must be a JSON object")
	}

	for _, field := range fields {
		if _, present := object[field]; !present {
			return refuse(fieldPath(path, field), "is missing")
		}
	}

	for field := range object {
		if !slices.Contains(fields, field) {
			return refuse(fieldPath(path, field), "is not a field of the schema")
		}
	}

	err := json.Unmarshal(raw, into)
	if err == nil {
		return nil
	}

	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) && typeError.Field != "" {
		return refuse(fieldPath(path, typeError.Field), "must be a "+typeError.Type.String()+", not a "+typeError.Value)
	}

	return refuse(objectName(path), "has a value of the wrong type: "+err.Error())
}
