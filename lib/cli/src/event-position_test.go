package main

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// eventPosition is where an observed line, by id, printed in its session,
// or -1 when it did not print.
func eventPosition(events []session.Event, id string) int {
	for position, event := range events {
		if event.ID == id {
			return position
		}
	}

	return -1
}
