package execution

import "errors"

// ErrCleanupUnconfirmed says nobody confirmed the runtime's cleanup: the
// outcome is unknown, the reservation is held, and the run stops.
var ErrCleanupUnconfirmed = errors.New("the runtime's cleanup was not confirmed")
