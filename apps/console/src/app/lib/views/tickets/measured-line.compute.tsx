import type { ReactNode } from 'react'

import type { TicketAnalysis } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'
import { UnknownValue } from './unknown-value'

// What the episodes measured: the resources, none cited once the episodes
// are judged, or a dash before that.
export function measuredLine(
	analysis: TicketAnalysis,
	format: TicketFormat,
): ReactNode {
	const { translate } = format
	if (analysis.measuredResources) {
		return (
			<span className="font-mono text-[12px]">
				{format.figures(analysis.measuredResources)}
			</span>
		)
	}
	if (!analysis.episodes) {
		return <UnknownValue className="font-mono text-[12px]" />
	}

	return (
		<>
			<span className="font-mono text-[12px]">
				{translate('analysis.measuredNone')}
			</span>
			<span className="text-[11px] text-muted-foreground">
				{translate('analysis.measuredReasoning')}
			</span>
		</>
	)
}
