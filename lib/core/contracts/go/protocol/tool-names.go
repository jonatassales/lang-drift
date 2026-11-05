package protocol

// The tools the server exposes to the agent.
const (
	// ToolProjectSync lists a project's runtimes without granting a start.
	ToolProjectSync = "langdrift.project.sync"

	// ToolProjectBind issues one environment credential for one runtime.
	ToolProjectBind = "langdrift.project.bind"

	// ToolBind binds one runtime with a runtime secret.
	ToolBind = "langdrift.bind"

	// ToolSync exchanges events and observations for the current intent.
	ToolSync = "langdrift.sync"

	// ToolClaim claims a start for one intent version and revision.
	ToolClaim = "langdrift.claim"
)
