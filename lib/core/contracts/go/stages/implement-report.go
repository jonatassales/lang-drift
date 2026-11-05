package stages

// ImplementReport is the implement skill's report: "tdd · 5 components · 14
// files · boundary only · stops at the reserved part".
type ImplementReport struct {
	Method                string `json:"method"`
	Components            int    `json:"components"`
	Files                 int    `json:"files"`
	BoundaryOnly          bool   `json:"boundaryOnly"`
	StoppedAtReservedPart bool   `json:"stoppedAtReservedPart"`
}
