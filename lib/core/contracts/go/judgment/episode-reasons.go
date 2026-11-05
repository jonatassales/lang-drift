package judgment

// The verdicts and the closed set of reasons an episode may get, so that a
// kept or dropped episode reads the same on every surface.
const (
	Kept    = "kept"
	Dropped = "dropped"

	// KeptReason is the one reason an episode is kept.
	KeptReason = "shares the stack and scale"
)

// DroppedReasons are the reasons an episode may be dropped.
var DroppedReasons = []string{"other runtime", "no broker", "docs only"}
