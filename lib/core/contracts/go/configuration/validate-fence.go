package configuration

// Validate returns every problem of a version 3 fence, never stopping at the
// first: its own fields, each component on its own, the paths, ports and
// images the components repeat, and the shares against the bound. The plan
// prints how many there are.
func (fence Fence) Validate() []Problem {
	problems := headerProblems(fence)

	for index, component := range fence.Components {
		problems = append(problems, componentProblems(index, component)...)
	}

	problems = append(problems, repeatedProblems(fence.Components)...)

	return append(problems, shareProblems(fence)...)
}
