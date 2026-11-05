package main

import (
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	clisession "github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

// lineTones are the escape sequences a terminal line of each kind opens with.
var lineTones = map[session.Kind]string{
	session.Comment: "\x1b[2m",
	session.Success: "\x1b[32m",
	session.Warning: "\x1b[33m",
}

// tonedLine is how a terminal shows one observed event: its words in its
// kind's tone, or plain when the kind keeps the terminal's color.
func tonedLine(event session.Event) string {
	words := clisession.Text(event)

	tone, colored := lineTones[clisession.KindOf(event)]
	if !colored {
		return words + "\n"
	}

	return tone + words + "\x1b[0m\n"
}
