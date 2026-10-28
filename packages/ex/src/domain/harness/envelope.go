package harness

// Envelope is the hardware envelope mirror: the ceiling one execution may
// take and the default a new runtime starts from.
type Envelope struct {
	Mirror           Mirror    `json:"mirror"`
	Version          int       `json:"version"`
	Status           string    `json:"status"`
	Owner            string    `json:"owner"`
	PerExecution     Resources `json:"perExecution"`
	Default          Resources `json:"default"`
	NotRuntimeFields []string  `json:"notRuntimeFields"`
}
