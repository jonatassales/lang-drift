package configuration_test

import (
	"encoding/json"
	"slices"
)

// sameKeys reports whether an object has exactly the named fields.
func sameKeys(names []string, object map[string]json.RawMessage) bool {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}

	slices.Sort(keys)
	sorted := slices.Sorted(slices.Values(names))

	return slices.Equal(keys, sorted)
}
