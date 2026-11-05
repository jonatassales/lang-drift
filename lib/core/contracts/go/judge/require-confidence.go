package judge

// requireConfidence refuses a confidence outside 0 to 1.
func requireConfidence(path string, value float64) error {
	if value < 0 || value > 1 {
		return refuse(path, "must be between 0 and 1")
	}

	return nil
}
