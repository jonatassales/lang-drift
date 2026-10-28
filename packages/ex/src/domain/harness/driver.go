package harness

// Driver is the coding agent the run drives for a stage's skills, as the
// owner configured it: an executable and its arguments, with placeholders
// substituted per argument, the environment variables passed through, a
// timeout, where usage is read in the result, and what to do when a call
// reports no usage.
type Driver struct {
	Executable     string   `json:"executable"`
	Args           []string `json:"args"`
	PassEnv        []string `json:"passEnv"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
	Usage          Usage    `json:"usage"`
	Unmeasured     string   `json:"unmeasured"`
}

// Usage names the fields of the driver's result that carry what a call used,
// as dotted paths into its JSON output; an empty path reads nothing.
type Usage struct {
	InputTokens  string `json:"inputTokens"`
	OutputTokens string `json:"outputTokens"`
	CostCents    string `json:"costCents"`
}
