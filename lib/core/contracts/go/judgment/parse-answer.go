package judgment

import "fmt"

// ParseAnswer reads the text of a judgment as the answer of its kind,
// refusing any field the kind does not define and any answer that breaks the
// kind's rules, naming the field and the rule.
func ParseAnswer(kind Kind, text string) (Answer, error) {
	answer := Answer{Kind: kind}

	var err error

	switch kind {
	case KindCriteria:
		answer.Criteria, err = parsed(text, checkCriteria)
	case KindEpisodes:
		answer.Episodes, err = parsed(text, checkEpisodes)
	case KindSuggestion:
		answer.Suggestion, err = parsed(text, checkSuggestion)
	default:
		return Answer{}, fmt.Errorf("judgment: no answer kind %q", kind)
	}

	if err != nil {
		return Answer{}, fmt.Errorf("judgment: the %s answer is refused: %w", kind, err)
	}

	return answer, nil
}
