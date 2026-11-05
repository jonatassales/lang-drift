package judge

import "strconv"

// The sets a reserved part draws from: the categories a part may belong to
// and the kinds the policy never reserves.
var (
	categories    = []string{"core-domain", "hard-part", "new-stack"}
	excludedKinds = []string{"boilerplate", "configuration", "generated", "rename", "tool-migration"}
)

// parseReservedPart reads a reserved-part answer. Empty paths mean no part
// qualifies and the reason says why; non-empty paths come with a goal,
// acceptance criteria and a category. The planned lines never exceed the
// ticket's total.
func parseReservedPart(text []byte) (*ReservedPartAnswer, error) {
	var answer ReservedPartAnswer
	if err := decodeObject(text, "", []string{"paths", "goal", "acceptance", "category", "plannedLines", "totalPlannedLines", "excludedKind", "reason", "confidence"}, &answer); err != nil {
		return nil, err
	}

	for index, path := range answer.Paths {
		if err := requireText("paths["+strconv.Itoa(index)+"]", path); err != nil {
			return nil, err
		}
	}
	for index, criterion := range answer.Acceptance {
		if err := requireText("acceptance["+strconv.Itoa(index)+"]", criterion); err != nil {
			return nil, err
		}
	}
	if answer.Category != nil {
		if err := requireOneOf("category", *answer.Category, categories); err != nil {
			return nil, err
		}
	}
	if answer.ExcludedKind != nil {
		if err := requireOneOf("excludedKind", *answer.ExcludedKind, excludedKinds); err != nil {
			return nil, err
		}
	}
	if answer.PlannedLines < 0 {
		return nil, refuse("plannedLines", "must not be negative")
	}
	if answer.TotalPlannedLines < 0 {
		return nil, refuse("totalPlannedLines", "must not be negative")
	}
	if answer.PlannedLines > answer.TotalPlannedLines {
		return nil, refuse("plannedLines", "must not exceed totalPlannedLines")
	}
	if err := requireText("reason", answer.Reason); err != nil {
		return nil, err
	}
	if err := requireConfidence("confidence", answer.Confidence); err != nil {
		return nil, err
	}

	proposed := len(answer.Paths) > 0
	if !proposed {
		return &answer, nil
	}

	if err := requireText("goal", answer.Goal); err != nil {
		return nil, err
	}
	if len(answer.Acceptance) == 0 {
		return nil, refuse("acceptance", "must name at least one criterion when paths are proposed")
	}
	if answer.Category == nil {
		return nil, refuse("category", "must be set when paths are proposed")
	}

	return &answer, nil
}
