package plan

// AskStatus is where one question of the planner stands. Before it, the
// money gates clear or pause it; after it, the judge route has answered,
// refused the secret, been unavailable, not offered the planner's model, not
// answered in time, relayed the model's refusal, or answered outside the
// contract.
type AskStatus string

// The statuses of a question, before and after it is asked.
const (
	AskCleared      AskStatus = "cleared"
	AskPaused       AskStatus = "paused"
	AskAnswered     AskStatus = "answered"
	AskUnauthorized AskStatus = "unauthorized"
	AskUnavailable  AskStatus = "unavailable"
	AskNotOffered   AskStatus = "not-offered"
	AskTimedOut     AskStatus = "timed-out"
	AskRefused      AskStatus = "refused"
	AskMalformed    AskStatus = "malformed"
)
