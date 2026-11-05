import type { ReactNode } from 'react'

import { revealed } from '@routes/platform/microruntimes/content/revealed.compute'
import { suggestedTotals } from '@routes/platform/microruntimes/content/suggested-totals.compute'

import { CardValue } from './card-value'
import { resolvedComponents } from './resolved-components.compute'
import type {
	TicketCardCopy,
	TicketCardFormat,
} from './ticket-card-format.compute'

export function planSummary(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	playedCount: number,
): ReactNode {
	const resolved = resolvedComponents(playedCount)

	if (resolved.length === 0) return <CardValue labels={copy} state="unset" />

	const figures = format.figures(suggestedTotals(resolved))

	if (!revealed(playedCount, 'drafted')) return figures

	return format.message('plan.summary', {
		figures,
		pending: copy.plan.pendingAuthorization,
	})
}
