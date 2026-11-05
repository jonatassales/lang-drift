package judge_test

import "slices"

// keys lists, sorted, the fields of a JSON object.
func keys(object map[string]any) []string {
	fields := make([]string, 0, len(object))
	for field := range object {
		fields = append(fields, field)
	}

	slices.Sort(fields)

	return fields
}
