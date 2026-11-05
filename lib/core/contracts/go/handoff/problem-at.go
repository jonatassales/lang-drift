package handoff

// problemAt is the problem of one field breaking one rule.
func problemAt(field string, rule string) Problem {
	return Problem{
		Field: field,
		Rule:  rule,
	}
}
