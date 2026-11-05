import {
	budgetCents,
	workflowDelivery,
	workflowFigures,
} from '@routes/platform/microruntimes/content/microruntimes-workflow'

import { runEnd } from './run-end.compute'
import type { RunOutcome } from './run-outcome.compute'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

export function workflowStatus(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	outcome: RunOutcome | null,
) {
	if (!outcome) return ''

	const end = runEnd(outcome)
	const figures = workflowFigures(outcome.choices)

	if (end === 'paused') {
		return format.message('workflow.statusPaused', {
			budget: format.money(budgetCents),
			reserved: format.money(figures.reserved),
		})
	}
	if (!end) return ''

	return format.message('workflow.statusDone', {
		number: workflowDelivery.pullRequest,
		observed: format.money(figures.observed),
		reserved: format.money(figures.reserved),
		state: copy.workflow[end],
	})
}
