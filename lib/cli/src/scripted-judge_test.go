package main

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// scriptedJudge answers each kind of question with a prepared text, the usage
// and the console's price of one call, or with the outcome it is told to
// give and the text of that outcome's kind; it remembers the connection and
// the questions it was asked.
type scriptedJudge struct {
	answers     map[string]string
	outcome     plan.AskOutcome
	asked       []plan.Question
	connections []plan.Connection
}

func (scripted *scriptedJudge) Ask(_ context.Context, connection plan.Connection, question plan.Question) (plan.Answer, plan.AskOutcome) {
	scripted.asked = append(scripted.asked, question)
	scripted.connections = append(scripted.connections, connection)

	if scripted.outcome.Status != "" {
		return plan.Answer{Text: scripted.answers[scripted.outcome.Kind]}, scripted.outcome
	}

	answer := plan.Answer{
		Text:         scripted.answers[question.Kind],
		InputTokens:  15000,
		OutputTokens: 1500,
		Cents:        5,
		Provenance:   "console",
		Illustrative: true,
		JudgeVersion: "console-judge-1",
	}
	outcome := plan.AskOutcome{
		Kind:   question.Kind,
		Status: plan.AskAnswered,
	}

	return answer, outcome
}
