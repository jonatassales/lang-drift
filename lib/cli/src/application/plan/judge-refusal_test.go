package plan_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

const consoleSecret = "judge-refusal-test-secret"

var recorded = plan.Connection{
	Recorded: true,
	URL:      "https://console.test",
	Secret:   consoleSecret,
}

// statuses are every status a question can reach, and one the planner never
// sets, which must refuse as well.
var statuses = []plan.AskStatus{
	plan.AskCleared,
	plan.AskPaused,
	plan.AskAnswered,
	plan.AskUnauthorized,
	plan.AskUnavailable,
	plan.AskNotOffered,
	plan.AskTimedOut,
	plan.AskRefused,
	plan.AskMalformed,
	plan.AskStatus("unknown"),
}

// Without a recorded connection there is no judge, so every question is
// refused, whatever its outcome, and the person is told to connect.
func TestJudgeRefusalWithoutAConnection(t *testing.T) {
	t.Parallel()

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()

			outcome := plan.AskOutcome{
				Kind:   plan.SourcesQuestion,
				Status: status,
				Reason: "the gate's or the console's words",
			}

			err := plan.JudgeRefusal(plan.Connection{}, outcome)

			if !errors.Is(err, plan.ErrNotConnected) || !strings.Contains(err.Error(), "dft agent connect") {
				t.Errorf("err = %v, want the person told to connect", err)
			}
		})
	}
}

// With a recorded connection, only a question the gates cleared, before it,
// and an answered question, after it, let the plan go on.
func TestJudgeRefusalLetsTheClearedAndTheAnsweredContinue(t *testing.T) {
	t.Parallel()

	for _, status := range []plan.AskStatus{plan.AskCleared, plan.AskAnswered} {
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()

			outcome := plan.AskOutcome{
				Kind:   plan.SuggestionQuestion,
				Status: status,
			}

			if err := plan.JudgeRefusal(recorded, outcome); err != nil {
				t.Errorf("err = %v, want the plan to continue", err)
			}
		})
	}
}

// Every other outcome refuses the plan with the cause in the gate's, the
// console's or the model's words, and says what the person does next; a
// gate's verdict carries its own next step.
func TestJudgeRefusalNamesTheCauseAndTheNextStep(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		status   plan.AskStatus
		reason   string
		sentinel error
		next     string
	}{
		{plan.AskPaused, "FG-04 reached · $0.16 observed of a $0.15 reservation; the calls before it stay on the ledger; raise the budget", plan.ErrPaused, "raise the budget"},
		{plan.AskUnauthorized, "the console refused RUNTIME_SECRET", plan.ErrModelUnavailable, "connect again"},
		{plan.AskUnavailable, "the console's judge is unavailable · platform-key-missing", plan.ErrModelUnavailable, "retry later"},
		{plan.AskNotOffered, "the console does not offer claude-sonnet-5-5; it offers claude-opus-5-5", plan.ErrModelUnavailable, "dft harness model set"},
		{plan.AskTimedOut, "the console did not answer in time", plan.ErrModelUnavailable, "retry later"},
		{plan.AskRefused, "the model stopped with refusal, not end_turn", plan.ErrJudgeRefused, "read the model's refusal"},
		{plan.AskMalformed, "judge: the sources answer is refused: answer: must be a JSON object", plan.ErrJudgeRefused, "retry later"},
		{plan.AskStatus("unknown"), "", plan.ErrJudgeRefused, "retry later"},
	} {
		t.Run(string(scenario.status), func(t *testing.T) {
			t.Parallel()

			outcome := plan.AskOutcome{
				Kind:   plan.SourcesQuestion,
				Status: scenario.status,
				Reason: scenario.reason,
			}

			err := plan.JudgeRefusal(recorded, outcome)

			if !errors.Is(err, scenario.sentinel) {
				t.Fatalf("err = %v, want %v", err, scenario.sentinel)
			}

			message := err.Error()
			if !strings.Contains(message, scenario.reason) || !strings.Contains(message, scenario.next) {
				t.Errorf("message = %q, want the cause %q and the next step %q", message, scenario.reason, scenario.next)
			}

			if strings.Contains(message, consoleSecret) {
				t.Errorf("the message carries the secret: %q", message)
			}
		})
	}
}

// A gate pause keeps the planner's message: the gate's detail, that nothing
// was paid, and the budget to raise, as the gate's verdict words it.
func TestJudgeRefusalKeepsTheGatePauseMessage(t *testing.T) {
	t.Parallel()

	outcome := plan.AskOutcome{
		Kind:   plan.SourcesQuestion,
		Status: plan.AskPaused,
		Reason: "FG-03 paused · a $0.15 reservation exceeds the $0.10 allocation; nothing was paid; raise the budget",
	}

	err := plan.JudgeRefusal(recorded, outcome)

	want := "the planner's paid calls are paused at a gate: FG-03 paused · a $0.15 reservation exceeds the $0.10 allocation; nothing was paid; raise the budget"
	if err == nil || err.Error() != want {
		t.Errorf("err = %v, want %q", err, want)
	}
}
