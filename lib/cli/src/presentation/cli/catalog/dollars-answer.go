package catalog

import (
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// dollarsAnswer reads an answer in dollars as cents, refusing it as invalid
// usage.
func dollarsAnswer(answer string) (int, error) {
	cents, err := finops.ParseDollars(answer)
	if err != nil {
		return 0, support.UsageError{Cause: err}
	}

	return cents, nil
}
