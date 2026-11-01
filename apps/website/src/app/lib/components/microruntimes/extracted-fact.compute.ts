import { type DocumentId, ticketPlan } from '@data/microruntimes-content'
import { budgetCents } from '@data/microruntimes-workflow'

import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

export function extractedFact(
	id: DocumentId,
	copy: TicketCardCopy,
	format: TicketCardFormat,
) {
	switch (id) {
		case 'prd':
			return copy.analysis.extracted.prd
		case 'designDoc':
			return format.message('analysis.extracted.designDoc', {
				budget: format.money(budgetCents),
				cpu: format.formatNumber.format(ticketPlan.bounded.cpuMillis),
				memory: format.formatNumber.format(ticketPlan.bounded.memoryMiB),
			})
		case 'spec':
			return format.message('analysis.extracted.spec', {
				components: ticketPlan.spec.components,
				criteria: ticketPlan.spec.acceptanceCriteria,
			})
		case 'guardrail':
			return format.message('analysis.extracted.guardrail', {
				criteria: ticketPlan.guardrailCriteria.join(' · '),
			})
	}
}
