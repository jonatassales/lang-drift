package session

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// Tee hands every event to the printed session and to an observer, in that
// order, and names the opened project to whichever of them wants it.
type Tee struct {
	Printed  session.Sink
	Observer session.Sink
}

// Emit hands the event to both.
func (tee Tee) Emit(event session.Event) {
	tee.Printed.Emit(event)
	tee.Observer.Emit(event)
}

// Opened names the project to both, when they want it.
func (tee Tee) Opened(project string) {
	session.Announce(tee.Printed, project)
	session.Announce(tee.Observer, project)
}
