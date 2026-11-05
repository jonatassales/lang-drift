package stages

import (
	"path"
	"strings"
)

// SpecReport is the to-spec skill's report: "0011-SPEC revised · 3
// acceptance criteria". Spec holds the spec's project-relative path.
type SpecReport struct {
	Spec               string `json:"spec"`
	Revised            bool   `json:"revised"`
	AcceptanceCriteria int    `json:"acceptanceCriteria"`
}

// ShortID is the spec's id as the stage line prints it: the first two
// dash-separated parts of its file name, 0011-SPEC for
// .agents/0011-SPEC-book-a-demo-runtime.spec.md. A file name with fewer
// parts is its own id.
func (report SpecReport) ShortID() string {
	parts := strings.SplitN(path.Base(report.Spec), "-", 3)
	if len(parts) < 2 {
		return path.Base(report.Spec)
	}

	return parts[0] + "-" + parts[1]
}
