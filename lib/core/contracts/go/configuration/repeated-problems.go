package configuration

// repeatedProblems are the components that repeat an earlier one's path,
// port or image: each of the three is unique across the runtime, so a later
// repeat is the problem.
func repeatedProblems(components []FenceComponent) []Problem {
	paths := map[string]bool{}
	ports := map[int]bool{}
	images := map[string]bool{}

	var problems []Problem

	for index, component := range components {
		if paths[component.Path] {
			problems = append(problems, problemAt(componentField(index, "path"), "must be unique"))
		}

		if ports[component.Port] {
			problems = append(problems, problemAt(componentField(index, "port"), "must be distinct"))
		}

		if images[component.Image] {
			problems = append(problems, problemAt(componentField(index, "image"), "must be distinct"))
		}

		paths[component.Path] = true
		ports[component.Port] = true
		images[component.Image] = true
	}

	return problems
}
