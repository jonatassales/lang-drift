import type { ReactNode } from 'react'

import type { TicketWorkflow } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'
import { workflowSummary } from './ticket-phases'

// The workflow's line: the model, effort and estimate against the budget
// while it runs, or how it ended.
export function workflowLine(
	workflow: TicketWorkflow,
	format: TicketFormat,
): ReactNode {
	const { translate } = format
	const summary = workflowSummary(workflow)
	switch (summary.kind) {
		case 'paused':
			return translate('workflow.summaryPaused')
		case 'done':
			return translate('workflow.summaryDone', {
				number: summary.number,
				observed: format.money(summary.observedCents),
				state: translate(`workflow.${summary.state}`),
			})
		case 'estimate':
			return translate('workflow.summaryChoice', {
				budget: format.money(workflow.budgetCents),
				effort: translate(`analysis.efforts.${workflow.effort}`),
				estimate: format.money(workflow.estimateCents),
				model: workflow.model.shortName,
			})
	}
}
