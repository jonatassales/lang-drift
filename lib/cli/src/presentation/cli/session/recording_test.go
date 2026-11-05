package session_test

import app "github.com/jonatassales/tao/lib/cli/src/application/session"

// recording keeps every event it receives, as a test's observer does.
type recording struct {
	events []app.Event
}

// Emit keeps the event.
func (sink *recording) Emit(event app.Event) {
	sink.events = append(sink.events, event)
}
