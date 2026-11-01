import type { ReactNode } from 'react'

import type { TicketPlan } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'
import { UnknownValue } from './unknown-value'

// The plan's line: a dash until it is drafted with totals, then the totals
// and whether the team has authorized them.
export function planLine(plan: TicketPlan, format: TicketFormat): ReactNode {
	if (!plan.suggested) return <UnknownValue />

	return format.t('plan.summary', {
		figures: format.figures(plan.suggested),
		pending: format.authorization(plan),
	})
}
