package catalog

import (
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/application/setup"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

	"github.com/spf13/pflag"
)

// initFlags are the answers ex init takes on its command line.
type initFlags struct {
	budget        string
	perRun        string
	reserveMargin int
	defaultGates  bool
	artifacts     string
}

// initAnswers reads the flags' answers. The three budget flags go together,
// and an amount that is not dollars is refused before anything is recorded.
func initAnswers(flags *pflag.FlagSet, values initFlags) (setup.Answers, error) {
	answers := setup.Answers{
		DefaultGates: values.defaultGates,
		Artifacts:    values.artifacts,
	}

	given := 0

	for _, name := range []string{"budget", "per-run", "reserve-margin"} {
		if flags.Changed(name) {
			given++
		}
	}

	if given == 0 {
		return answers, nil
	}

	if given < 3 {
		return setup.Answers{}, support.UsageError{Cause: fmt.Errorf("--budget, --per-run and --reserve-margin go together")}
	}

	periodCents, err := finops.ParseDollars(values.budget)
	if err != nil {
		return setup.Answers{}, support.UsageError{Cause: fmt.Errorf("--budget: %w", err)}
	}

	perRunCents, err := finops.ParseDollars(values.perRun)
	if err != nil {
		return setup.Answers{}, support.UsageError{Cause: fmt.Errorf("--per-run: %w", err)}
	}

	answers.Budget = &finops.BudgetBounds{
		PeriodCents:          periodCents,
		PerRunCents:          perRunCents,
		ReserveMarginPercent: values.reserveMargin,
	}

	return answers, nil
}
