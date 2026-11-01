import type { TicketAnalysis, TicketDocumentId } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'

// What the read of each document extracted, or null while it is unread or
// the fact is not known yet.
export function extractedFact(
	id: TicketDocumentId,
	analysis: TicketAnalysis,
	format: TicketFormat,
): string | null {
	if (!analysis.documentsRead.includes(id)) return null

	const { number, translate } = format

	switch (id) {
		case 'prd':
			return analysis.useCase
				? translate('analysis.extracted.prd', { useCase: analysis.useCase })
				: null

		case 'designDoc':
			if (!analysis.bounded || analysis.budgetCents === null) return null

			return translate('analysis.extracted.designDoc', {
				budget: format.money(analysis.budgetCents),
				cpu: number.format(analysis.bounded.cpuMillis),
				memory: number.format(analysis.bounded.memoryMiB),
			})

		case 'spec':
			if (
				analysis.specComponents === null ||
				analysis.acceptanceCriteria === null
			)
				return null

			return translate('analysis.extracted.spec', {
				components: analysis.specComponents,
				criteria: analysis.acceptanceCriteria,
			})

		case 'guardrail':
			if (!analysis.criteria) return null

			return translate('analysis.extracted.guardrail', {
				criteria: analysis.criteria.join(' · '),
			})
	}
}
