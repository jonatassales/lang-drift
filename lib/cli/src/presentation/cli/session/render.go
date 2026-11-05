package session

import (
	"encoding/json"
	"io"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// Human prints each line as the website's terminal shows it. On a terminal
// each line carries its kind's tone and the window's title names the
// project; elsewhere the bytes are the plain text.
type Human struct {
	Writer   io.Writer
	Terminal bool
}

// Emit prints one line.
func (human Human) Emit(event session.Event) {
	line := Text(event)
	if human.Terminal {
		line = toned(KindOf(event), line)
	}

	writeTo(human.Writer, line+"\n")
}

// Opened titles a terminal's window with the tool and the project's
// folder, "dft — project".
func (human Human) Opened(project string) {
	if !human.Terminal {
		return
	}

	writeTo(human.Writer, "\x1b]0;dft — "+project+"\x07")
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

// For picks the renderer the global flags ask for, toned when the output is
// a terminal.
func For(writer io.Writer, asJSON bool, terminal bool) session.Sink {
	if asJSON {
		return NDJSON{Writer: writer}
	}

	return Human{
		Writer:   writer,
		Terminal: terminal,
	}
}
