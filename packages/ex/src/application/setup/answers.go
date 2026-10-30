// Package setup completes ex init: after the project identity, it records the
// inputs a team gives on the way in, through the same use cases as their own
// commands, and never replaces one already recorded.
package setup

import "github.com/jonatassales/tao/packages/ex/src/domain/finops"

// Answers are what the flags and the questions gave. A nil budget is no
// answer.
type Answers struct {
	Budget       *finops.BudgetBounds
	DefaultGates bool
}
