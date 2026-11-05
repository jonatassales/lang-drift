package catalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jonatassales/tao/lib/cli/src/application/setup"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// initGuide asks ex init's questions on the terminal, for the pending inputs
// the flags left unanswered.
type initGuide struct {
	prompter *support.Prompter
}

// Ask asks for the month's budget, which an empty answer skips, then for the
// default gates, the artifacts folder and CI Gates.
func (guide initGuide) Ask(ctx context.Context, pending setup.Pending, answers setup.Answers) (setup.Answers, error) {
	if pending.Budget && answers.Budget == nil {
		budget, err := guide.budget(ctx, pending.Period)
		if err != nil {
			return setup.Answers{}, err
		}

		answers.Budget = budget
	}

	if pending.Gates && !answers.DefaultGates {
		apply, err := guide.prompter.Confirm(ctx, "Apply the default gates, FG-01 to FG-06? [y/N]")
		if err != nil {
			return setup.Answers{}, err
		}

		answers.DefaultGates = apply
	}

	if pending.Artifacts && answers.Artifacts == "" {
		folder, err := guide.prompter.Optional(ctx, "Artifacts folder, where the tickets and specs live, relative to the project (empty skips it)")
		if err != nil {
			return setup.Answers{}, err
		}

		answers.Artifacts = folder
	}

	if pending.CIGates && !answers.CIGates {
		generate, err := guide.prompter.Confirm(ctx, "Generate the CI Gates workflow, .github/workflows/langdrift-gates.yml? [y/N]")
		if err != nil {
			return setup.Answers{}, err
		}

		answers.CIGates = generate
	}

	return answers, nil
}

// budget asks for the period's amount, the per-run cap and the reserve
// margin. An empty amount skips the budget; the margin has no default,
// because the limit is the team's to choose.
func (guide initGuide) budget(ctx context.Context, period string) (*finops.BudgetBounds, error) {
	amount, err := guide.prompter.Optional(ctx, "Monthly budget for "+period+", in dollars (empty skips the budget)")
	if err != nil || amount == "" {
		return nil, err
	}

	periodCents, err := dollarsAnswer(amount)
	if err != nil {
		return nil, err
	}

	perRun, err := guide.prompter.Optional(ctx, "Per-run cap, in dollars")
	if err != nil {
		return nil, err
	}

	perRunCents, err := dollarsAnswer(perRun)
	if err != nil {
		return nil, err
	}

	margin, err := guide.prompter.Optional(ctx, "Reserve margin, in percent")
	if err != nil {
		return nil, err
	}

	marginPercent, err := strconv.Atoi(margin)
	if err != nil {
		return nil, support.UsageError{Cause: fmt.Errorf("the reserve margin %q is not a whole percentage", margin)}
	}

	return &finops.BudgetBounds{
		PeriodCents:          periodCents,
		PerRunCents:          perRunCents,
		ReserveMarginPercent: marginPercent,
	}, nil
}
