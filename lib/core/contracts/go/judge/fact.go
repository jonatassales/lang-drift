package judge

// Fact is one statement the material makes, with its citation.
type Fact struct {
	Topic  string `json:"topic"`
	Text   string `json:"text"`
	Source string `json:"source"`
}
