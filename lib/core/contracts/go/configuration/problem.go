package configuration

// Problem is one rule a version 3 fence breaks: the field, as a dotted path
// such as components[1].port, and the rule in a few words.
type Problem struct {
	Field string
	Rule  string
}

// String reads the problem as the field and its rule.
func (problem Problem) String() string {
	return problem.Field + " " + problem.Rule
}
