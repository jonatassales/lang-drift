package gates

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/domain/gates"
)

// ReportRequest names the repository and the pull request's range.
type ReportRequest struct {
	Directory string
	Base      string
	Head      string
}

// ReportResult is what the report shows: the range, its commits and the
// FinOps readings.
type ReportResult struct {
	Base    string          `json:"base"`
	Head    string          `json:"head"`
	Commits []gates.Commit  `json:"commits"`
	FinOps  []gates.Reading `json:"finops"`
}

// Report reads a pull request's range for CI Gates.
type Report struct {
	History History
}

// Execute lists the range's commits and leaves every FinOps reading unknown.
func (report Report) Execute(ctx context.Context, request ReportRequest) (ReportResult, error) {
	commits, err := report.History.Commits(ctx, request.Directory, request.Base, request.Head)
	if err != nil {
		return ReportResult{}, err
	}

	return ReportResult{
		Base:    request.Base,
		Head:    request.Head,
		Commits: commits,
		FinOps:  gates.UnknownFinOps(),
	}, nil
}
