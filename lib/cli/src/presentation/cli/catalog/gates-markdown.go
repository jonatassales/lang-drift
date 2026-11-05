package catalog

import (
	"fmt"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application/gates"
)

// gatesMarkdown renders the CI Gates report for a job summary: the range and
// its commits, then the FinOps readings.
func gatesMarkdown(result gates.ReportResult) string {
	var report strings.Builder

	report.WriteString("## LangDrift CI Gates\n\n")

	noun := "commits"
	if len(result.Commits) == 1 {
		noun = "commit"
	}

	fmt.Fprintf(&report, "Range `%s..%s`: %d %s.\n\n", result.Base, result.Head, len(result.Commits), noun)

	if len(result.Commits) > 0 {
		report.WriteString("| Commit | Subject |\n| --- | --- |\n")

		for _, commit := range result.Commits {
			fmt.Fprintf(&report, "| `%s` | %s |\n", commit.SHA, markdownCell(commit.Subject))
		}

		report.WriteString("\n")
	}

	report.WriteString("### FinOps\n\n| Reading | Status | Reason |\n| --- | --- | --- |\n")

	for _, reading := range result.FinOps {
		fmt.Fprintf(&report, "| %s | %s | %s |\n", reading.Name, reading.Status, markdownCell(reading.Reason))
	}

	return report.String()
}
