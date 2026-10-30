package catalog

import (
	"context"
	"fmt"
	"os"

	"github.com/jonatassales/tao/packages/ex/src/application/gates"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// GatesReporter makes the CI Gates report of a pull request.
type GatesReporter interface {
	Execute(context.Context, gates.ReportRequest) (gates.ReportResult, error)
}

// NewGatesReportCommand creates the CI Gates report command. It needs no
// project: BASE_SHA and HEAD_SHA name the pull request's range, and the report
// goes to the file GITHUB_STEP_SUMMARY names, else to standard output.
func NewGatesReportCommand(reporter GatesReporter, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Report a pull request's commit range for CI Gates; FinOps stays unknown without a committed source",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			base := os.Getenv("BASE_SHA")
			if base == "" {
				return support.UsageError{Cause: fmt.Errorf("BASE_SHA must name the pull request's base commit")}
			}

			head := os.Getenv("HEAD_SHA")
			if head == "" {
				return support.UsageError{Cause: fmt.Errorf("HEAD_SHA must name the pull request's head commit")}
			}

			ctx, cancel := globals.Bounded(command.Context())
			defer cancel()

			result, err := reporter.Execute(ctx, gates.ReportRequest{
				Directory: globals.Fallback,
				Base:      base,
				Head:      head,
			})
			if err != nil {
				return err
			}

			if globals.JSON {
				return support.Render(command.OutOrStdout(), result, true)
			}

			summary := os.Getenv("GITHUB_STEP_SUMMARY")
			if summary == "" {
				_, err = fmt.Fprint(command.OutOrStdout(), gatesMarkdown(result))

				return err
			}

			return appendJobSummary(summary, gatesMarkdown(result))
		},
	}
}
