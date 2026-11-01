import type { ReactNode } from 'react'

import type { TicketAnalysis } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'
import { analysisSummary } from './ticket-phases'
import { UnknownValue } from './unknown-value'

// The analysis's line: the latest thing it has settled.
export function analysisLine(
	analysis: TicketAnalysis,
	format: TicketFormat,
): ReactNode {
	const { translate } = format
	const summary = analysisSummary(analysis)
	switch (summary.kind) {
		case 'suggestion':
			return translate('analysis.suggestionSummary', {
				effort: translate(`analysis.efforts.${summary.effort}`),
				model: summary.model,
			})
		case 'episodes':
			return translate('analysis.episodesSummary', {
				kept: summary.kept,
				total: summary.total,
			})
		case 'documents':
			return translate('analysis.documentsSummary', { count: summary.count })
		case 'unknown':
			return <UnknownValue />
	}
}
