package judge

// FieldError says which field of an answer broke which rule, as
// ParseAnswer reports a refusal.
type FieldError struct {
	Field string
	Rule  string
}

func (failure FieldError) Error() string {
	return failure.Field + ": " + failure.Rule
}
