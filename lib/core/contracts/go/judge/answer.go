package judge

import "fmt"

// ParseAnswer turns the text of a judgment into the typed answer of its
// kind, reading it as the kind's schema does, and refuses anything else with
// a FieldError that names the field and the rule it broke. The console has
// validated the text once; the planner validates it again here before it
// trusts a single field.
func ParseAnswer(kind Kind, text string) (Answer, error) {
	answer := Answer{Kind: kind}
	raw := []byte(text)

	var err error
	switch kind {
	case KindSources:
		answer.Sources, err = parseSources(raw)
	case KindEpisodes:
		answer.Episodes, err = parseEpisodes(raw)
	case KindSuggestion:
		answer.Suggestion, err = parseSuggestion(raw)
	case KindReservedPart:
		answer.ReservedPart, err = parseReservedPart(raw)
	default:
		return Answer{}, fmt.Errorf("judge: no answer kind %q", kind)
	}

	if err != nil {
		return Answer{}, fmt.Errorf("judge: the %s answer is refused: %w", kind, err)
	}

	return answer, nil
}
