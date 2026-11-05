package domain

import "strings"

// reservedCategories are the categories a reserved part may belong to.
var reservedCategories = []string{"core-domain", "hard-part", "new-stack"}

// generateReservedPart reserves the rule the ticket turns on, in its first
// component, with its tests. Every fifth part is the lead's pick; the rest
// are the judge's.
func generateReservedPart(number int, work workTemplate) *TicketReservedPart {
	component := work.components[0]
	name := strings.ReplaceAll(work.useCase, " ", "-")

	source := component.Path + "/src/domain/" + name + ".ts"
	test := component.Path + "/src/domain/" + name + ".test.ts"
	if component.Technology == "Go" {
		source = component.Path + "/domain/" + name + ".go"
		test = component.Path + "/domain/" + name + "_test.go"
	}

	chosenBy := "judge"
	if number%5 == 0 {
		chosenBy = "lead"
	}

	return &TicketReservedPart{
		Paths:    []string{source, test},
		Goal:     "Write the " + work.useCase + " rule the ticket turns on, with its tests.",
		Category: reservedCategories[number%len(reservedCategories)],
		ChosenBy: chosenBy,
	}
}
