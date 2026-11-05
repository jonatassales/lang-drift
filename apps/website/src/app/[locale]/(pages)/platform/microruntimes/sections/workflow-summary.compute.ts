import {
	budgetCents,
	type WorkflowChoices,
	workflowDelivery,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import { workflowFigures } from '@routes/platform/microruntimes/content/workflow-figures.compute'
import { workflowModel } from '@routes/platform/microruntimes/content/workflow-model.compute'

import { runEnd } from './run-end.compute'
import type { RunOutcome } from './run-outcome.compute'
import type {
	TicketCardCopy,
	TicketCardFormat,
} from './ticket-card-format.factory'

export function workflowSummary(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	choices: WorkflowChoices,
	outcome: RunOutcome | null,
) {
	const end = outcome ? runEnd(outcome) : null

	if (end === 'paused') return copy.workflow.summaryPaused
	if (outcome && end) {
		return format.message('workflow.summaryDone', {
			number: workflowDelivery.pullRequest,
			observed: format.money(workflowFigures(outcome.choices).observed),
			state: copy.workflow[end],
		})
	}

	return format.message('workflow.summaryChoice', {
		budget: format.money(budgetCents),
		effort: copy.analysis.efforts[choices.effort],
		estimate: format.money(workflowFigures(choices).estimate),
		model: workflowModel(choices.model).shortName,
	})
}
