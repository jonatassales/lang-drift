package judge

// fieldPath names a field inside its parent, such as sources[0].reason; a
// field of the answer itself is named alone.
func fieldPath(parent string, field string) string {
	if parent == "" {
		return field
	}

	return parent + "." + field
}
