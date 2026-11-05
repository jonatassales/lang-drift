import {
	budgetCents,
	expectedTokens,
	observedTokens,
	type WorkflowChoices,
} from './microruntimes-workflow.data'
import { reserveCents } from './reserve-cents.compute'
import { spendCents } from './spend-cents.compute'

// The run's money in cents: the estimate, the reservation the gate would
// hold, what the usage record observed and what the reservation releases.
export function workflowFigures(choices: WorkflowChoices) {
	const estimate = spendCents(choices.model, expectedTokens[choices.effort])
	const reserved = reserveCents(estimate)
	const observed = spendCents(choices.model, observedTokens[choices.effort])

	return {
		estimate,
		reserved,
		observed,
		released: reserved - observed,
		overBudget: reserved > budgetCents,
	}
}

export type WorkflowFigures = ReturnType<typeof workflowFigures>
