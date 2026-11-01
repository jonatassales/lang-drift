import type { TicketPlan } from '@data/types'

import type { PlanFigures } from './ticket-phases.types.ts'

// The plan's totals, what they leave in the envelope and the share they
// take of it. Without suggested totals there is no headroom, never a zero
// one, and the bars stay empty.
export function planFigures(plan: TicketPlan): PlanFigures {
	const { bounded, suggested } = plan
	if (!suggested) {
		return {
			suggested: null,
			headroom: null,
			cpuShare: 0,
			memoryShare: 0,
		}
	}

	return {
		suggested,
		headroom: {
			cpuMillis: bounded.cpuMillis - suggested.cpuMillis,
			memoryMiB: bounded.memoryMiB - suggested.memoryMiB,
		},
		cpuShare: suggested.cpuMillis / bounded.cpuMillis,
		memoryShare: suggested.memoryMiB / bounded.memoryMiB,
	}
}
