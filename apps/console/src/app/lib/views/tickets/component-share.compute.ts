import type { Resources, TicketComponent } from '@data/types'

// A component's share of the plan's CPU envelope, from 0 to 1.
export function componentShare(
	component: Pick<TicketComponent, 'cpuMillis'>,
	bounded: Resources,
) {
	return component.cpuMillis / bounded.cpuMillis
}
