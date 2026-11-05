package stages

// checked reads a report of one stage and applies the stage's rules.
func checked[Value any](stage Stage, text string, check func(Value) error) (*Value, error) {
	var value Value
	if err := readReport(text, reportFields[stage], &value); err != nil {
		return nil, err
	}

	if err := check(value); err != nil {
		return nil, err
	}

	return &value, nil
}
