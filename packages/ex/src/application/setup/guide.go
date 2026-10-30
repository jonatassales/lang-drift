package setup

import "context"

// Guide asks a person for the pending inputs the flags left unanswered.
type Guide interface {
	Ask(ctx context.Context, pending Pending, answers Answers) (Answers, error)
}
