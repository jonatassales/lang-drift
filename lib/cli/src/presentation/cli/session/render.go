package session

import (
	"encoding/json"
	"io"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// Human prints each line as the website's terminal shows it.
type Human struct {
	Writer io.Writer
}

// Emit prints one line.
func (human Human) Emit(event session.Event) {
	_, _ = io.WriteString(human.Writer, Text(event)+"\n")
}

// NDJSON prints one JSON event per line, with the website's line and reveal
// ids, so a console or a TUI fills the same cards in the same order.
type NDJSON struct {
	Writer io.Writer
	// Now stamps the events; nil uses the wall clock.
	Now func() time.Time
}

// Record is one event as NDJSON prints it.
type Record struct {
	At      string       `json:"at"`
	Session string       `json:"session"`
	Kind    session.Kind `json:"kind"`
	ID      string       `json:"id"`
	Text    string       `json:"text"`
	Card    string       `json:"card,omitempty"`
	Field   string       `json:"field,omitempty"`
	Data    session.Data `json:"data,omitempty"`
}

// Emit prints one event.
func (ndjson NDJSON) Emit(event session.Event) {
	record := Record{
		At:      ndjson.now().UTC().Format(time.RFC3339),
		Session: event.Session,
		Kind:    KindOf(event),
		ID:      event.ID,
		Text:    Text(event),
		Data:    event.Data,
	}

	if reveal, known := RevealOf(event); known {
		record.Card, record.Field = reveal.Card, reveal.Field
	}

	encoder := json.NewEncoder(ndjson.Writer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(record)
}

func (ndjson NDJSON) now() time.Time {
	if ndjson.Now != nil {
		return ndjson.Now()
	}

	return time.Now()
}

// For picks the renderer the global flags ask for.
func For(writer io.Writer, asJSON bool) session.Sink {
	if asJSON {
		return NDJSON{Writer: writer}
	}

	return Human{Writer: writer}
}
