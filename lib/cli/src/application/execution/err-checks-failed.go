package execution

import "errors"

// ErrChecksFailed says the pull request's checks ended with a failure, which
// the simulation never shows: the run stops once it has recorded and
// settled the lap.
var ErrChecksFailed = errors.New("the pull request's checks failed")
