package handoff

// Authorization is where the plan's authorization stands when the plan
// hands off.
type Authorization string

// AuthorizationPending is a plan the team has not authorized yet.
const AuthorizationPending Authorization = "pending"
