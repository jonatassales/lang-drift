import { hasPlayed, ticketPlan } from '@data/microruntimes-content'

export function resolvedComponents(playedCount: number) {
	return ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
}
