package judge

// requireText refuses an empty string where the schema wants at least one
// character.
func requireText(path string, value string) error {
	if value == "" {
		return refuse(path, "must not be empty")
	}

	return nil
}
