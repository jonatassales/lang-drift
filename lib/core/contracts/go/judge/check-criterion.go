package judge

// checkCriterion refuses a criterion without its id or what it applies to.
func checkCriterion(path string, criterion Criterion) error {
	if err := requireText(path+".id", criterion.ID); err != nil {
		return err
	}

	return requireText(path+".applies", criterion.Applies)
}
