package protocol

// ProjectMembership lists a project's runtimes. Discovery never grants a start.
type ProjectMembership struct {
	Repository Repository `json:"repository"`
	ProjectID  string     `json:"projectId"`
	RuntimeIDs []string   `json:"runtimeIds"`
}
