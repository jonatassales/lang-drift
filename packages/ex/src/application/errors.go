// Package application holds EX's use cases: what a user can ask EX to do, and
// the order in which the domain and the outside world are consulted to do it.
//
// It depends on the domain and on nothing concrete. Every capability it needs
// from the outside world is a small port declared in ports.go; infrastructure
// implements those ports and the composition root injects them. Each use case
// is a struct of the ports it needs, with one Execute method taking a request
// and returning a result. On error, the result is always the zero value.
package application

import (
	"context"
	"errors"
	"fmt"
)

// Conditions adapters report so use cases can act on them. Adapters build
// these errors with NotFound and AlreadyExists; callers match them with
// errors.Is.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// Application errors classify refusals that span several domain objects.
var (
	ErrAdoptionRefused      = errors.New("refusing to adopt what EX did not create")
	ErrTargetMismatch       = errors.New("target differs from the recorded one")
	ErrConfigurationChanged = errors.New("configuration changed since execution")
)

// NotFound reports that something an adapter was asked for does not exist.
func NotFound(format string, arguments ...any) error {
	return conditionError{condition: ErrNotFound, message: fmt.Sprintf(format, arguments...)}
}

// AlreadyExists reports that an adapter refused to replace something.
func AlreadyExists(format string, arguments ...any) error {
	return conditionError{condition: ErrAlreadyExists, message: fmt.Sprintf(format, arguments...)}
}

// conditionError carries a condition without repeating it in its message, so
// "no cluster is recorded" does not read as "not found: no cluster is recorded".
type conditionError struct {
	condition error
	message   string
}

func (err conditionError) Error() string { return err.message }

// Is lets errors.Is match the condition.
func (err conditionError) Is(target error) bool { return target == err.condition }

// isInterruption reports whether an operation stopped because it was canceled
// or ran out of time, rather than because it failed.
func isInterruption(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
