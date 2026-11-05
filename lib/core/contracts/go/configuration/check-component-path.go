package configuration

// CheckComponentPath accepts the repository root, ".", as a component source,
// and otherwise the same paths as CheckRelativePath.
func CheckComponentPath(text string) error {
	if text == "." {
		return nil
	}

	return CheckRelativePath(text)
}
