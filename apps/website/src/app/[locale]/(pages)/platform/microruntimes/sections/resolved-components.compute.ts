import { hasPlayed } from '@routes/platform/microruntimes/content/has-played.compute'
import { ticketPlan } from '@routes/platform/microruntimes/content/microruntimes-content.data'

export function resolvedComponents(playedCount: number) {
	return ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
}
