package handoff

// Delivery is the run's branch, commit and pull request (spec 0063). Its
// other fields join with that spec.
type Delivery struct {
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
}
