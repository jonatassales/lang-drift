package stages

// checkPlanning accepts a planning report with a count of settled questions.
func checkPlanning(report PlanningReport) error {
	return checkCount("openQuestionsSettled", report.OpenQuestionsSettled)
}
