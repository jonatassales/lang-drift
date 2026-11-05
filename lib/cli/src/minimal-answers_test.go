package main

import "github.com/jonatassales/tao/lib/cli/src/application/plan"

// minimalAnswers are the smallest judgments the contract accepts: no verdict,
// fact or episode, and the recorded model at medium effort.
func minimalAnswers() map[string]string {
	return map[string]string{
		plan.SourcesQuestion:    `{"sources":[],"facts":[],"criteria":[],"figures":[]}`,
		plan.EpisodesQuestion:   `{"episodes":[]}`,
		plan.SuggestionQuestion: `{"model":"claude-sonnet-5-5","effort":"medium","effortReason":"a bounded change","modelReason":"the recorded model fits","reasons":[]}`,
	}
}
