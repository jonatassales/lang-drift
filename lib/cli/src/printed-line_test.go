package main

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// printedLine is one line of an executed session: its line id, its kind and
// its printed text.
type printedLine struct {
	id   string
	kind session.Kind
	text string
}
