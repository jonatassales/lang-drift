package judge

// objectName names an object in a refusal: the answer itself, or the item
// its path points at.
func objectName(path string) string {
	if path == "" {
		return "answer"
	}

	return path
}
