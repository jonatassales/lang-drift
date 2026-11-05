package stages

// checkCount refuses a negative count, naming its field.
func checkCount(field string, value int) error {
	if value < 0 {
		return FieldError{
			Field: field,
			Rule:  "must not be negative",
		}
	}

	return nil
}
