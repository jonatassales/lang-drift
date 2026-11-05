package main

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// sessionRecorder observes a session as the dependencies' Observer port: it
// keeps every event, so a test reads each line's id, kind and words.
type sessionRecorder struct {
	events []session.Event
}

// Emit keeps the event.
func (recorder *sessionRecorder) Emit(event session.Event) {
	recorder.events = append(recorder.events, event)
}
