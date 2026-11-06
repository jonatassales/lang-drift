package execution

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// stopSilently records why the run stops, the first cause only, and drops
// every line after it: the runtime goes down without its line, and the run
// returns the cause once it has recorded and settled what happened.
func (playing *playing) stopSilently(cause error) {
	if playing.stopCause == nil {
		playing.stopCause = cause
	}

	playing.sink = session.Discard{}
}
