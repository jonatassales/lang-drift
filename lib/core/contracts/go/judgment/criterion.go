package judgment

// Criterion is one guardrail criterion the judge finds applies, as
// "FIN-02 estimate before a paid call" prints it.
type Criterion struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
}
