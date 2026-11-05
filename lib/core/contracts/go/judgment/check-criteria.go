package judgment

import "strconv"

// checkCriteria accepts a criteria answer whose every entry names its
// criterion once and says what it requires.
func checkCriteria(answer CriteriaAnswer) error {
	seen := map[string]bool{}

	for index, criterion := range answer.Criteria {
		field := "criteria[" + strconv.Itoa(index) + "]"

		namedOnce := criterion.ID != "" && !seen[criterion.ID]
		if !namedOnce {
			return FieldError{
				Field: field + ".id",
				Rule:  "must name a criterion once",
			}
		}

		if criterion.Statement == "" {
			return FieldError{
				Field: field + ".statement",
				Rule:  "must say what the criterion requires",
			}
		}

		seen[criterion.ID] = true
	}

	return nil
}
