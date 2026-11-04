package harness

// Workflow is the delivery-stages mirror: the team's stages, where a person
// takes part, and the efforts a run may choose from.
type Workflow struct {
	Mirror   Mirror   `json:"mirror"`
	Version  int      `json:"version"`
	Status   string   `json:"status"`
	Owner    string   `json:"owner"`
	Note     string   `json:"note"`
	Stages   []Stage  `json:"stages"`
	Defaults Defaults `json:"defaults"`
	Efforts  []string `json:"efforts"`
}

// Stage is one delivery stage: its skills and whether a person takes part,
// required, optional or owner-only.
type Stage struct {
	ID     string   `json:"id"`
	Label  string   `json:"label"`
	Skills []string `json:"skills"`
	Human  string   `json:"human"`
	// Hooks are committed scripts the stage runs first.
	Hooks []string `json:"hooks,omitempty"`
}

// Defaults names the stages a person takes part in unless a run says otherwise.
type Defaults struct {
	HumanAt []string `json:"humanAt"`
}

// How a person takes part in a stage.
const (
	HumanRequired = "required"
	HumanOptional = "optional"
	OwnerOnly     = "owner-only"
)
