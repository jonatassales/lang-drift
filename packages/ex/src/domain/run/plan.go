// Package run holds the records of one runtime's run: the plan the planner
// wrote, the execution record, and the JSON fences of the memory and the
// laps that the RAG indexes. Field order follows the records as the console
// reads them; a value nobody measured is null, never zero.
package run

import "github.com/jonatassales/tao/packages/ex/src/domain/harness"

// Plan is the analysis and the plan of one run: the ticket, what its links
// reached, the facts and inputs, the episodes retrieved, the suggestion, the
// components and their sizing, the checks and the authorized runtime file.
type Plan struct {
	Version         int             `json:"version"`
	Sample          bool            `json:"sample,omitempty"`
	Run             string          `json:"run"`
	State           string          `json:"state"`
	Ticket          Ticket          `json:"ticket"`
	Collection      Collection      `json:"collection"`
	Sources         []Source        `json:"sources"`
	Facts           []Fact          `json:"facts"`
	Inputs          Inputs          `json:"inputs"`
	Discrepancies   []Discrepancy   `json:"discrepancies"`
	Retrieval       Retrieval       `json:"retrieval"`
	Sizing          Sizing          `json:"sizing"`
	Suggestion      Suggestion      `json:"suggestion"`
	Components      []Component     `json:"components"`
	Totals          Totals          `json:"totals"`
	Materialization Materialization `json:"materialization"`
	Price           Price           `json:"price"`
	Checks          []Check         `json:"checks"`
	Runtime         RuntimeRef      `json:"runtime"`
	ContractGaps    []string        `json:"contractGaps"`
	Authorization   *Authorization  `json:"authorization,omitempty"`
}

// Authorization is who authorized the plan, and when: running it is the
// authorization.
type Authorization struct {
	By  string            `json:"by"`
	Via string            `json:"via"`
	At  harness.Timestamp `json:"at"`
}

// Ticket is the ticket the person named, read by the tracker's convention.
type Ticket struct {
	Path    string  `json:"path"`
	Number  *string `json:"number"`
	Slug    string  `json:"slug"`
	Title   *string `json:"title"`
	Status  *string `json:"status"`
	NamedBy Actor   `json:"namedBy"`
}

// Actor is who did something and through what.
type Actor struct {
	Person string `json:"person"`
	Via    string `json:"via"`
}

// Collection is how the sources were gathered: the agent that judged them,
// the limits, what the traversal reached and the context hash.
type Collection struct {
	Agent   *CollectionAgent `json:"agent"`
	Limits  Limits           `json:"limits"`
	Reached Reached          `json:"reached"`
	Context Context          `json:"context"`
}

// CollectionAgent is the model that kept or discarded each source; its calls
// are paid work on the ledger's planning line.
type CollectionAgent struct {
	Model         string `json:"model"`
	PromptVersion string `json:"promptVersion"`
	Paid          bool   `json:"paid"`
	LedgerLine    string `json:"ledgerLine"`
}

// Limits bound the traversal.
type Limits struct {
	MaxDepth       int    `json:"maxDepth"`
	MaxFiles       int    `json:"maxFiles"`
	MaxBytes       int64  `json:"maxBytes"`
	Secrets        string `json:"secrets"`
	OutsideProject string `json:"outsideProject"`
	Symlinks       string `json:"symlinks"`
}

// Reached is how far the traversal went.
type Reached struct {
	Files int   `json:"files"`
	Depth int   `json:"depth"`
	Bytes int64 `json:"bytes"`
}

// Context is the hash of everything read, and the rule that produced it.
type Context struct {
	SHA256 string `json:"sha256"`
	Rule   string `json:"rule"`
}

// Source is one file the ticket's links reached, with the agent's verdict.
type Source struct {
	Path         string   `json:"path"`
	MediaType    *string  `json:"mediaType"`
	Bytes        *int64   `json:"bytes"`
	SHA256       *string  `json:"sha256"`
	Depth        int      `json:"depth"`
	LinkedFrom   *string  `json:"linkedFrom"`
	LinkSyntax   string   `json:"linkSyntax"`
	State        string   `json:"state"`
	Verdict      *string  `json:"verdict"`
	Reason       *string  `json:"reason"`
	Confidence   *float64 `json:"confidence"`
	InferredRole *string  `json:"inferredRole"`
}

// Fact is something a source said, by topic, with its citation.
type Fact struct {
	Topic  string `json:"topic"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

// Inputs are the values that bound the run, each with the record that set it.
type Inputs struct {
	CPUMillis        InputValue `json:"cpuMillis"`
	MemoryMiB        InputValue `json:"memoryMiB"`
	ModelBudgetCents InputValue `json:"modelBudgetCents"`
	Ceiling          Ceiling    `json:"ceiling"`
	NotRuntimeFields []string   `json:"notRuntimeFields"`
}

// InputValue is one explicit input and its provenance.
type InputValue struct {
	Value  int               `json:"value"`
	Record string            `json:"record"`
	SetBy  string            `json:"setBy"`
	Via    string            `json:"via"`
	SetAt  harness.Timestamp `json:"setAt"`
}

// Ceiling is the envelope the inputs stayed inside.
type Ceiling struct {
	CPUMillis int               `json:"cpuMillis"`
	MemoryMiB int               `json:"memoryMiB"`
	Record    string            `json:"record"`
	SetBy     string            `json:"setBy"`
	Via       string            `json:"via"`
	SetAt     harness.Timestamp `json:"setAt"`
}

// Discrepancy is a fact compared with an input; the input holds either way.
type Discrepancy struct {
	Topic  string `json:"topic"`
	Fact   string `json:"fact"`
	Input  string `json:"input"`
	Result string `json:"result"`
}

// Retrieval is how earlier runs were searched and which were considered.
type Retrieval struct {
	Method        string       `json:"method"`
	Criteria      []string     `json:"criteria"`
	KeepAtMatches int          `json:"keepAtMatches"`
	Considered    []Considered `json:"considered"`
	Note          string       `json:"note"`
}

// Considered is one earlier run the retrieval scored.
type Considered struct {
	Run        string   `json:"run"`
	Similarity float64  `json:"similarity"`
	Criteria   []string `json:"criteria"`
	Score      float64  `json:"score"`
	Verdict    string   `json:"verdict"`
	Reason     string   `json:"reason"`
}

// Sizing says what the proposed shares were based on.
type Sizing struct {
	Basis             string   `json:"basis"`
	MeasuredRunsCited []string `json:"measuredRunsCited"`
	Note              string   `json:"note"`
}

// Suggestion is the model and effort recommended, and who authorized them.
type Suggestion struct {
	By           string   `json:"by"`
	Model        string   `json:"model"`
	ModelName    *string  `json:"modelName"`
	Effort       string   `json:"effort"`
	Reasons      []string `json:"reasons"`
	AuthorizedBy string   `json:"authorizedBy"`
}

// Component is one component of the runtime, as inferred, confirmed, proposed
// and executed.
type Component struct {
	Path         string            `json:"path"`
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	Technologies []Technology      `json:"technologies"`
	Dockerfile   string            `json:"dockerfile"`
	Image        string            `json:"image"`
	Port         int               `json:"port"`
	Readiness    *Readiness        `json:"readiness"`
	EnvFile      *string           `json:"envFile"`
	PassEnv      []string          `json:"passEnv"`
	InferredFrom []string          `json:"inferredFrom"`
	ConfirmedBy  Actor             `json:"confirmedBy"`
	Proposed     harness.Resources `json:"proposed"`
	Executed     Executed          `json:"executed"`
}

// Readiness is the probe a component's workload waits on, and where it
// came from: the Dockerfile's HEALTHCHECK, or the fence's explicit kind.
// Period, timeout and initial delay are seconds.
type Readiness struct {
	Kind                string   `json:"kind"`
	From                string   `json:"from"`
	Path                *string  `json:"path"`
	Command             []string `json:"command"`
	PeriodSeconds       int      `json:"periodSeconds"`
	TimeoutSeconds      int      `json:"timeoutSeconds"`
	FailureThreshold    int      `json:"failureThreshold"`
	InitialDelaySeconds int      `json:"initialDelaySeconds"`
}

// Technology is one technology a component runs on.
type Technology struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Executed is the share a component actually got, and the rule that gave it.
type Executed struct {
	CPUMillis int    `json:"cpuMillis"`
	MemoryMiB int    `json:"memoryMiB"`
	Rule      string `json:"rule"`
}

// Totals are the proposed shares against the envelope.
type Totals struct {
	Proposed harness.Resources `json:"proposed"`
	Bounded  harness.Resources `json:"bounded"`
	Headroom harness.Resources `json:"headroom"`
}

// Materialization is what the engine generates from the plan.
type Materialization struct {
	Deployments  int `json:"deployments"`
	ReplicasEach int `json:"replicasEach"`
}

// Price is the price basis the plan had.
type Price struct {
	Record        *string `json:"record"`
	Mirror        *string `json:"mirror"`
	ReadOn        *string `json:"readOn"`
	RuntimePriced bool    `json:"runtimePriced"`
	Note          string  `json:"note"`
}

// Check is one cross-check of the draft.
type Check struct {
	ID     string  `json:"id"`
	Result string  `json:"result"`
	Errors *int    `json:"errors,omitempty"`
	Detail *string `json:"detail,omitempty"`
}

// RuntimeRef names the runtime file the plan authorized.
type RuntimeRef struct {
	Path     string `json:"path"`
	Revision int    `json:"revision"`
	SHA256   string `json:"sha256"`
	State    string `json:"state"`
}

// SplitEnvelope divides an envelope equally among count components, the
// remainder going to the first ones: the share runtime.md v2 can express.
func SplitEnvelope(envelope harness.Resources, count int) []harness.Resources {
	if count <= 0 {
		return nil
	}

	shares := make([]harness.Resources, count)
	for index := range shares {
		shares[index] = harness.Resources{CPUMillis: envelope.CPUMillis / count, MemoryMiB: envelope.MemoryMiB / count}
		if index < envelope.CPUMillis%count {
			shares[index].CPUMillis++
		}
		if index < envelope.MemoryMiB%count {
			shares[index].MemoryMiB++
		}
	}

	return shares
}

// EqualSplitRule is how the executed shares were given.
const EqualSplitRule = "equal split of the runtime envelope; runtime.md v2 carries no per-component share"

// The contract gaps a plan records while the fence is at version 2.
var ContractGaps = []string{
	"per-component cpu and memory are not runtime.md v2 fields, so the proposal and the execution can differ",
	"technologies are not runtime.md v2 fields; the plan keeps them for display",
}
