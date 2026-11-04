// Package gates holds what CI Gates reports about a pull request: its commits
// and the readings it can, or cannot, make.
package gates

// Commit is one commit of a pull request's range.
type Commit struct {
	SHA     string `json:"sha"`
	Subject string `json:"subject"`
}
