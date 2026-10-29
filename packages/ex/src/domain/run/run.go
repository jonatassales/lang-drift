package run

import "github.com/jonatassales/tao/packages/ex/src/domain/harness"

// Run is the execution record of one runtime: its laps, its stages, its
// delivery and its outcome.
type Run struct {
	Version              int                   `json:"version"`
	Sample               bool                  `json:"sample,omitempty"`
	ID                   string                `json:"id"`
	Kind                 string                `json:"kind"`
	Ticket               *string               `json:"ticket"`
	Runtime              string                `json:"runtime"`
	State                string                `json:"state"`
	Outcome              *string               `json:"outcome"`
	ProposedTicketStatus *ProposedTicketStatus `json:"proposedTicketStatus"`
	Window               Window                `json:"window"`
	Files                Files                 `json:"files"`
	HumanInTheLoop       HumanInTheLoop        `json:"humanInTheLoop"`
	Laps                 []LapSummary          `json:"laps"`
	Stages               []StageState          `json:"stages"`
	Delivery             Delivery              `json:"delivery"`
}

// ProposedTicketStatus is the status the run proposes; the owner decides.
type ProposedTicketStatus struct {
	Status    string `json:"status"`
	DecidedBy string `json:"decidedBy"`
	Decided   bool   `json:"decided"`
}

// Window is when the run started and finished.
type Window struct {
	Started  *harness.Timestamp `json:"started"`
	Finished *harness.Timestamp `json:"finished"`
}

// Files names the run's records.
type Files struct {
	Plan     *string `json:"plan"`
	Memory   string  `json:"memory"`
	Ledger   string  `json:"ledger"`
	Evidence string  `json:"evidence"`
}

// HumanInTheLoop names the workflow record and the stages a person chose.
type HumanInTheLoop struct {
	Record *string  `json:"record"`
	Mirror *string  `json:"mirror"`
	Chosen []string `json:"chosen"`
}

// LapSummary is one lap as the run lists it.
type LapSummary struct {
	Number        string  `json:"number"`
	File          string  `json:"file"`
	Choice        Choice  `json:"choice"`
	Outcome       string  `json:"outcome"`
	Gate          *string `json:"gate,omitempty"`
	Execution     *string `json:"execution,omitempty"`
	NothingPaid   bool    `json:"nothingPaid"`
	EstimateCents int     `json:"estimateCents,omitempty"`
	ReservedCents int     `json:"reservedCents,omitempty"`
}

// Choice is the model and effort a lap ran with, and who chose them.
type Choice struct {
	Model    string  `json:"model"`
	Effort   string  `json:"effort"`
	By       string  `json:"by"`
	Overrode *string `json:"overrode"`
	Via      string  `json:"via"`
	Record   string  `json:"record"`
}

// StageState is where one workflow stage stands.
type StageState struct {
	ID          string             `json:"id"`
	State       string             `json:"state"`
	Human       bool               `json:"human"`
	ConfirmedBy *string            `json:"confirmedBy,omitempty"`
	At          *harness.Timestamp `json:"at"`
}

// Delivery is what the run delivered, as observed.
type Delivery struct {
	Branch      *string       `json:"branch"`
	Commit      *string       `json:"commit"`
	PullRequest *PullRequest  `json:"pullRequest"`
	Checks      []CheckResult `json:"checks"`
}

// PullRequest is the pull request the run observed.
type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
}

// CheckResult is one of the pull request's checks.
type CheckResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}
