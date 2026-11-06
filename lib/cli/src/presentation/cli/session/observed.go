package session

import "github.com/jonatassales/tao/lib/cli/src/application/session"

// Observed lets an observer, such as a test, see every event the printed
// session receives; without an observer the printed session stands alone.
func Observed(printed session.Sink, observer session.Sink) session.Sink {
	if observer == nil {
		return printed
	}

	return Tee{
		Printed:  printed,
		Observer: observer,
	}
}
