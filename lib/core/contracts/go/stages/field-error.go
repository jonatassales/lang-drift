package stages

// FieldError says which field of a report broke which rule.
type FieldError struct {
	Field string
	Rule  string
}

// Error reads the refusal as the field and its rule.
func (failure FieldError) Error() string {
	return failure.Field + ": " + failure.Rule
}
