import type { ReactNode } from 'react'

import type { TicketPlan } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'
import { UnknownValue } from './unknown-value'

// The line under the components: how many have resolved of the spec's
// count while the plan sizes, then the count and the totals.
export function componentsLine(
	plan: TicketPlan,
	specComponents: number | null,
	format: TicketFormat,
): ReactNode {
	const { number, t } = format
	const resolved = plan.components.reduce(
		(totals, component) => ({
			cpuMillis: totals.cpuMillis + component.cpuMillis,
			memoryMiB: totals.memoryMiB + component.memoryMiB,
		}),
		{
			cpuMillis: 0,
			memoryMiB: 0,
		},
	)
	const values = {
		cpu: number.format(resolved.cpuMillis),
		memory: number.format(resolved.memoryMiB),
	}
	if (plan.suggested) {
		return t('plan.componentsTotal', {
			count: plan.components.length,
			...values,
		})
	}
	if (specComponents === null) return <UnknownValue />

	return t('plan.resolvedProgress', {
		resolved: plan.components.length,
		total: specComponents,
		...values,
	})
}
