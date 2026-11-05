package judge_test

// requestSchemaProperty is the part of a property's schema the request uses.
type requestSchemaProperty struct {
	Type      string   `json:"type"`
	Enum      []string `json:"enum"`
	MinLength *int     `json:"minLength"`
	Minimum   *float64 `json:"minimum"`
	Maximum   *float64 `json:"maximum"`
}
