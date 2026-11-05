package judge

// factTopics are the topics a fact may speak of.
var factTopics = []string{"resource", "component", "acceptance", "budget", "limits"}

// checkFact refuses a fact with a topic outside the set, or without its
// text or its source.
func checkFact(path string, fact Fact) error {
	if err := requireOneOf(path+".topic", fact.Topic, factTopics); err != nil {
		return err
	}
	if err := requireText(path+".text", fact.Text); err != nil {
		return err
	}

	return requireText(path+".source", fact.Source)
}
