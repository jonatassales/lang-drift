package execution

import "errors"

// ErrPausedAtGate says FG-03 paused the run before anything was paid: the
// reservation exceeds the budget. It is a shown ending, "# needs a decision:
// raise the budget or lower the effort · nothing was paid for", not a stop.
var ErrPausedAtGate = errors.New("the reservation exceeds the budget")
