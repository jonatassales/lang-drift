package handoff

// Human says, per stage, whether a person takes part in it: the HUMAN IN THE
// LOOP ticks of the Workflow card.
type Human struct {
	Planning  bool `json:"planning"`
	Spec      bool `json:"spec"`
	Tickets   bool `json:"tickets"`
	Implement bool `json:"implement"`
	Review    bool `json:"review"`
}
