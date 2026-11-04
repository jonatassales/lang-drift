package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// maxNesting bounds how deeply a consumer-authored JSON document may nest.
const maxNesting = 64

// decodeStrict rejects duplicate keys, unknown struct fields, and trailing
// documents, then decodes data into target. Data must already be bounded.
func decodeStrict(data []byte, target any) error {
	scanner := json.NewDecoder(bytes.NewReader(data))

	if err := scanValue(scanner, 0); err != nil {
		return err
	}

	if _, err := scanner.Token(); !errors.Is(err, io.EOF) {
		return errors.New("expected exactly one JSON document")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	return decoder.Decode(target)
}

func scanValue(scanner *json.Decoder, depth int) error {
	if depth > maxNesting {
		return fmt.Errorf("JSON nesting exceeds %d levels", maxNesting)
	}

	token, err := scanner.Token()
	if err != nil {
		return err
	}

	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return nil
	}

	switch delimiter {
	case '{':
		if err := scanObject(scanner, depth); err != nil {
			return err
		}

	case '[':
		for scanner.More() {
			if err := scanValue(scanner, depth+1); err != nil {
				return err
			}
		}

	default:
		return errors.New("unexpected JSON delimiter")
	}

	_, err = scanner.Token()

	return err
}

func scanObject(scanner *json.Decoder, depth int) error {
	keys := map[string]bool{}

	for scanner.More() {
		token, err := scanner.Token()
		if err != nil {
			return err
		}

		key, isKey := token.(string)
		if !isKey {
			return errors.New("expected a JSON object key")
		}

		if keys[key] {
			return fmt.Errorf("duplicate JSON key %q", key)
		}

		keys[key] = true

		if err := scanValue(scanner, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// requireDescriptorFields rejects missing or null required descriptor fields,
// which a struct decoder would silently read as zero values.
func requireDescriptorFields(data []byte) error {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}

	if err := requireKeys(root, []string{"version", "id", "environment", "owner", "units", "terraform", "kubernetes"}, "hooks"); err != nil {
		return err
	}

	if raw, present := root["hooks"]; present && isNull(raw) {
		return errors.New("hooks must be an object")
	}

	if err := requireNested(root["terraform"], "directory"); err != nil {
		return fmt.Errorf("terraform: %w", err)
	}

	if err := requireNested(root["kubernetes"], "context", "namespace"); err != nil {
		return fmt.Errorf("kubernetes: %w", err)
	}

	var units map[string]map[string]json.RawMessage
	if err := json.Unmarshal(root["units"], &units); err != nil {
		return err
	}

	for name, unit := range units {
		if err := requireKeys(unit, []string{"component", "kubernetes"}); err != nil {
			return fmt.Errorf("unit %s: %w", name, err)
		}

		if err := requireNested(unit["kubernetes"], "manifests"); err != nil {
			return fmt.Errorf("unit %s: %w", name, err)
		}
	}

	return nil
}

func requireNested(raw json.RawMessage, required ...string) error {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return err
	}

	return requireKeys(object, required)
}

// requireKeys requires every required key to be present and not null, and
// rejects keys that are neither required nor optional.
func requireKeys(object map[string]json.RawMessage, required []string, optional ...string) error {
	allowed := map[string]bool{}

	for _, key := range required {
		allowed[key] = true

		if value, present := object[key]; !present || isNull(value) {
			return fmt.Errorf("required field %q is missing or null", key)
		}
	}

	for _, key := range optional {
		allowed[key] = true
	}

	for key := range object {
		if !allowed[key] {
			return fmt.Errorf("unknown field %q", key)
		}
	}

	return nil
}

func isNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}
