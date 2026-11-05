package plan

import "fmt"

// JudgeRefusal decides whether the plan goes on: before a question, with the
// money gates' verdict, and after it, with the console's outcome. Only a
// recorded connection with a question the gates cleared or the judge answered
// continues, so no path reaches a proposed plan without judged answers. Every
// other case refuses the plan with the cause in the gate's, the console's or
// the model's words and what the person does next: connect, choose an offered
// model, retry later, or read the model's refusal; a gate's verdict carries
// its own next step. A status the planner never sets refuses too.
func JudgeRefusal(connection Connection, outcome AskOutcome) error {
	if !connection.Recorded {
		return fmt.Errorf("%w: connect it with dft agent connect <console-url>, then plan again", ErrNotConnected)
	}

	switch outcome.Status {
	case AskCleared, AskAnswered:
		return nil

	case AskPaused:
		return fmt.Errorf("%w: %s", ErrPaused, outcome.Reason)

	case AskUnauthorized:
		return fmt.Errorf("%w: %s; connect again with the console's current RUNTIME_SECRET", ErrModelUnavailable, outcome.Reason)

	case AskNotOffered:
		return fmt.Errorf("%w: %s; choose one of them with dft harness model set <model> --planner <offered model>", ErrModelUnavailable, outcome.Reason)

	case AskUnavailable, AskTimedOut:
		return fmt.Errorf("%w: %s; retry later", ErrModelUnavailable, outcome.Reason)

	case AskRefused:
		return fmt.Errorf("%w: the model refused the %s question · %s; read the model's refusal before planning again", ErrJudgeRefused, outcome.Kind, outcome.Reason)
	}

	return fmt.Errorf("%w: the %s answer is outside the contract · %s; retry later", ErrJudgeRefused, outcome.Kind, outcome.Reason)
}
