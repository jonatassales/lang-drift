package judge_test

import "slices"

// propertyViolations checks one value against its property's schema.
func propertyViolations(field string, property requestSchemaProperty, value any) []string {
	text, isText := value.(string)
	number, isNumber := value.(float64)
	isInteger := isNumber && number == float64(int64(number))

	switch {
	case len(property.Enum) > 0 && !slices.Contains(property.Enum, text):
		return []string{field + " is outside its enum"}
	case property.Type == "string" && !isText:
		return []string{field + " is not a string"}
	case property.Type == "integer" && !isInteger:
		return []string{field + " is not an integer"}
	case property.MinLength != nil && len(text) < *property.MinLength:
		return []string{field + " is shorter than its minimum"}
	case property.Minimum != nil && number < *property.Minimum:
		return []string{field + " is below its minimum"}
	case property.Maximum != nil && number > *property.Maximum:
		return []string{field + " is above its maximum"}
	}

	return nil
}
