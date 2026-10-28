package harness

// Choice is the default model and effort a run takes, and the model the
// planner's own calls use. It is CLI input held privately; no mirror shows it,
// because it blocks nothing by itself.
type Choice struct {
	Model   string `json:"model"`
	Effort  string `json:"effort"`
	Planner string `json:"planner"`
}

// Efforts are the efforts a run may choose from.
var Efforts = []string{"low", "medium", "high"}

// ValidEffort reports whether text names one of the efforts.
func ValidEffort(text string) bool {
	for _, effort := range Efforts {
		if text == effort {
			return true
		}
	}

	return false
}
