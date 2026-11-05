package judgment_test

import (
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/judgment"
)

// Each kind's sample is the answer the simulation shows, and parses into its
// type: the criteria of "console › judging · FIN-02 estimate before a paid
// call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to
// the load", the six EPISODES rows, and the suggestion behind "✓ suggests
// claude-sonnet-5-5 · effort medium".
func TestParseAnswerReadsTheSimulationsAnswers(t *testing.T) {
	t.Parallel()

	criteria, err := judgment.ParseAnswer(judgment.KindCriteria, readSample(t, "criteria.answer.sample.json"))
	if err != nil || len(criteria.Criteria.Criteria) != 3 || criteria.Criteria.Criteria[0].Statement != "estimate before a paid call" {
		t.Errorf("criteria = %+v, %v", criteria.Criteria, err)
	}

	episodes, err := judgment.ParseAnswer(judgment.KindEpisodes, readSample(t, "episodes.answer.sample.json"))
	if err != nil || len(episodes.Episodes.Episodes) != 6 || episodes.Episodes.Episodes[3].Reason != "docs only" {
		t.Errorf("episodes = %+v, %v", episodes.Episodes, err)
	}

	suggestion, err := judgment.ParseAnswer(judgment.KindSuggestion, readSample(t, "suggestion.answer.sample.json"))
	suggested := err == nil && suggestion.Suggestion.Model == "claude-sonnet-5-5" && suggestion.Suggestion.Effort == "medium"
	if !suggested || !strings.HasPrefix(suggestion.Suggestion.EffortReason, "low skips") {
		t.Errorf("suggestion = %+v, %v", suggestion.Suggestion, err)
	}
}

// An answer outside its kind is refused naming the field: a field the kind
// does not define, a price the plan never pays, a criterion without its
// statement or repeated, a suggestion without its effort or reasons, and an
// unknown kind.
func TestParseAnswerRefusesWhatTheKindDoesNotDefine(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		kind  judgment.Kind
		text  string
		field string
	}{
		{judgment.KindCriteria, `{"criteria":[{"id":"FIN-02","statement":"estimate before a paid call","applies":"x"}]}`, "answer"},
		{judgment.KindCriteria, `{"criteria":[{"id":"FIN-02","statement":""}]}`, "criteria[0].statement"},
		{judgment.KindCriteria, `{"criteria":[{"id":"FIN-02","statement":"a"},{"id":"FIN-02","statement":"b"}]}`, "criteria[1].id"},
		{judgment.KindSuggestion, `{"model":"claude-sonnet-5-5","effort":"max","effortReason":"a","modelReason":"b"}`, "effort"},
		{judgment.KindSuggestion, `{"model":"claude-sonnet-5-5","effort":"medium","effortReason":"","modelReason":"b"}`, "effortReason"},
		{judgment.KindSuggestion, `{"model":"","effort":"medium","effortReason":"a","modelReason":"b"}`, "model"},
		{judgment.KindSuggestion, `{"model":"claude-sonnet-5-5","effort":"medium","effortReason":"a","modelReason":"b","cost":{"cents":1}}`, "answer"},
		{judgment.KindEpisodes, `{"episodes":[],"usage":{"inputTokens":1}}`, "answer"},
	} {
		_, err := judgment.ParseAnswer(scenario.kind, scenario.text)
		if err == nil || !strings.Contains(err.Error(), scenario.field+":") {
			t.Errorf("%s %s gave %v, want a refusal naming %s", scenario.kind, scenario.text, err, scenario.field)
		}
	}

	if _, err := judgment.ParseAnswer("sources", `{}`); err == nil {
		t.Error("the sources kind, which the plan no longer asks, was accepted")
	}
}
