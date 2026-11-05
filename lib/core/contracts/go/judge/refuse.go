package judge

// refuse names the field of an answer that broke a rule of its schema.
func refuse(field string, rule string) FieldError {
	return FieldError{
		Field: field,
		Rule:  rule,
	}
}
