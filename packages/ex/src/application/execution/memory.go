package execution

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// memory composes what the run teaches: the laps' story, the lessons the
// figures hold, and the fence the RAG indexes.
func (playing *playing) memory(outcome string) Memory {
	lessons, facts := playing.lessons()

	memory := Memory{
		Runtime:        playing.runtime.String(),
		Name:           playing.config.Name,
		Ticket:         path.Base(playing.planned.Ticket.Path),
		Kind:           playing.record.Kind,
		Outcome:        playing.outcomeWords(outcome),
		ProposedStatus: "none proposed",
		Story:          playing.story(),
		Lessons:        lessons,
		Fence: run.RunFence{
			Run:               playing.runtime.String(),
			Ticket:            text(ticketReference(playing.planned)),
			Kind:              playing.record.Kind,
			Outcome:           text(fenceOutcome(outcome)),
			Laps:              len(playing.record.Laps),
			ComponentPaths:    playing.componentPaths(),
			Stacks:            playing.stacks(),
			MeasuredResources: playing.lap.Resources != nil && playing.lap.Resources.CPUPeak != nil && playing.lap.Resources.MemoryPeak != nil,
			MeasuredTokens:    !playing.lap.Unmeasured && (playing.lap.InputTokens > 0 || playing.lap.OutputTokens > 0),
			Observed:          playing.observed(),
			Facts:             facts,
		},
	}

	if proposed := playing.record.ProposedTicketStatus; proposed != nil {
		memory.ProposedStatus = proposed.Status + ", pending the " + proposed.DecidedBy
	}

	return memory
}

func fenceOutcome(outcome string) string {
	if outcome == Unknown {
		return Unknown
	}

	return Completed
}

func (playing *playing) outcomeWords(outcome string) string {
	laps := strconv.Itoa(len(playing.record.Laps)) + " laps"
	if len(playing.record.Laps) == 1 {
		laps = "one lap"
	}
	switch outcome {
	case Unknown:
		return "unknown after " + laps + ": the runtime's cleanup was not confirmed"
	case ReadyToMerge:
		return "completed in " + laps + ", ready to merge"
	}

	return "completed in " + laps + ", ready for review"
}

// story is one sentence per lap.
func (playing *playing) story() string {
	sentences := make([]string, 0, len(playing.record.Laps))
	for _, lap := range playing.record.Laps {
		sentences = append(sentences, playing.lapSentence(lap))
	}

	return strings.Join(sentences, " ")
}

func (playing *playing) lapSentence(lap run.LapSummary) string {
	choice := playing.modelName(lap.Choice.Model) + " at " + lap.Choice.Effort + " effort"
	if lap.Choice.Overrode != nil {
		choice = "the owner's choice over the " + *lap.Choice.Overrode + ", " + choice
	} else {
		choice = "the suggestion, " + choice
	}

	switch {
	case lap.Outcome == PausedAtGate && lap.NothingPaid:
		return "Lap " + lap.Number + " stopped in the pit: " + choice + " asked for a reservation above the per-run cap, and nothing was paid."
	case lap.Outcome == PausedAtGate:
		return "Lap " + lap.Number + " paused its paid calls at FG-04 with " + finops.Dollars(lap.ObservedCents) + " observed."
	case lap.Outcome == BlockedAt:
		return "Lap " + lap.Number + " was blocked at the gate before anything was paid."
	case lap.Outcome == Completed:
		return "Lap " + lap.Number + " followed " + choice + " and finished with " + finops.Dollars(lap.ObservedCents) + " observed of " + finops.Dollars(lap.ReservedCents) + " reserved."
	case lap.Outcome == Unknown:
		return "Lap " + lap.Number + " ran " + choice + "; its outcome is unknown because the runtime's cleanup was not confirmed."
	}

	return "Lap " + lap.Number + " ran " + choice + "."
}

// lessons are what the figures teach, as prose and as the slugs the RAG
// indexes.
func (playing *playing) lessons() ([]string, []string) {
	var lessons, facts []string
	cap := playing.guards.Budget.Bounded.PerRunCents

	for _, lap := range playing.record.Laps {
		if lap.Outcome != PausedAtGate || !lap.NothingPaid {
			continue
		}
		lessons = append(lessons, "At this ticket's size, "+playing.modelName(lap.Choice.Model)+" at "+lap.Choice.Effort+" effort needs a reservation above the "+finops.Dollars(cap)+" cap: "+finops.Dollars(lap.ReservedCents)+" was asked for.")
		facts = append(facts, slug(lap.Choice.Model)+"-"+lap.Choice.Effort+"-exceeds-per-run-cap")
	}

	if resources := playing.lap.Resources; resources != nil && resources.CPUPeak != nil && resources.MemoryPeak != nil {
		lessons = append(lessons, fmt.Sprintf("%d MiB of %s MiB and %dm of %dm at peak: the next plan can propose with this headroom in mind.", *resources.MemoryPeak, finops.Grouped(resources.Executed.MemoryMiB), *resources.CPUPeak, resources.Executed.CPUMillis))
		facts = append(facts, fmt.Sprintf("memory-peak-%d-of-%d-mib", *resources.MemoryPeak, resources.Executed.MemoryMiB), fmt.Sprintf("cpu-peak-%d-of-%dm", *resources.CPUPeak, resources.Executed.CPUMillis))
	}

	if observed := playing.observed(); observed.TokensThousands != nil {
		tokens := *observed.TokensThousands
		lessons = append(lessons, playing.modelName(playing.lap.Choice.Model)+" at "+playing.lap.Choice.Effort+" effort used "+finops.Tokens(tokens.Input)+" input and "+finops.Tokens(tokens.Output)+" output tokens; the next estimate for this ticket can start from these instead of the effort table.")
		facts = append(facts, slug(playing.lap.Choice.Model)+"-"+playing.lap.Choice.Effort+"-used-"+strings.ToLower(finops.Tokens(tokens.Input))+"-input-"+strings.ToLower(finops.Tokens(tokens.Output))+"-output")
	}

	if playing.lap.Ready != nil {
		for _, component := range playing.config.Components {
			lessons = append(lessons, "`"+component.ReadinessPath+"` passed readiness for "+component.Path+".")
			facts = append(facts, "readiness-"+slug(component.ReadinessPath))
		}
	}

	if checks := playing.record.Delivery.Checks; len(checks) > 0 && !anyPending(checks) {
		passed := 0
		for _, check := range checks {
			if check.Result == CheckPassed {
				passed++
			}
		}
		facts = append(facts, fmt.Sprintf("checks-%d-of-%d-passed", passed, len(checks)))
	}

	if lessons == nil {
		lessons = []string{}
	}
	if facts == nil {
		facts = []string{}
	}

	return lessons, facts
}

func (playing *playing) modelName(model string) string {
	if price, found := playing.guards.Price(model); found {
		return price.Name
	}

	return model
}

func (playing *playing) componentPaths() []string {
	paths := []string{}
	for _, component := range playing.planned.Components {
		paths = append(paths, component.Path)
	}

	return paths
}

func (playing *playing) stacks() []string {
	stacks := []string{}
	seen := map[string]bool{}
	for _, component := range playing.planned.Components {
		for _, technology := range component.Technologies {
			if seen[technology.ID] {
				continue
			}
			seen[technology.ID] = true
			stacks = append(stacks, technology.ID)
		}
	}

	return stacks
}

// slug makes a lowercase, dash-separated word of a value.
func slug(value string) string {
	var out strings.Builder
	dash := false
	for _, character := range strings.ToLower(value) {
		isWord := character >= 'a' && character <= 'z' || character >= '0' && character <= '9'
		switch {
		case isWord:
			out.WriteRune(character)
			dash = false
		case !dash && out.Len() > 0:
			out.WriteByte('-')
			dash = true
		}
	}

	return strings.TrimSuffix(out.String(), "-")
}
