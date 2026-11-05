package handoff

// Problem is one rule a hand-off record breaks: the field and the rule in a
// few words. It is the hand-off's own, apart from the fence's, so that the
// two contracts change on their own.
type Problem struct {
	Field string
	Rule  string
}

// String reads the problem as the field and its rule.
func (problem Problem) String() string {
	return problem.Field + " " + problem.Rule
}
