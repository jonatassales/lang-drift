package judgment

// parsed decodes an answer strictly into its type and applies its check.
func parsed[Value any](text string, check func(Value) error) (*Value, error) {
	var value Value
	if err := decodeStrictly(text, &value); err != nil {
		return nil, err
	}

	if err := check(value); err != nil {
		return nil, err
	}

	return &value, nil
}
