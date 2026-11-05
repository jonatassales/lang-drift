package stages

// reportFields are the fields each stage's report holds, every one required.
var reportFields = map[Stage][]string{
	Planning:  {"openQuestionsSettled"},
	Spec:      {"spec", "revised", "acceptanceCriteria"},
	Tickets:   {"tickets", "blockingEdgesRecorded"},
	Implement: {"method", "components", "files", "boundaryOnly", "stoppedAtReservedPart"},
	Review:    {"findings", "fixed"},
}
