package judge_test

// requestSchema is the contract's request.schema.json, as the console reads it.
type requestSchema struct {
	Required             []string                         `json:"required"`
	AdditionalProperties bool                             `json:"additionalProperties"`
	Properties           map[string]requestSchemaProperty `json:"properties"`
}
