package stages

// checkReview accepts a review report with a count of findings.
func checkReview(report ReviewReport) error {
	return checkCount("findings", report.Findings)
}
