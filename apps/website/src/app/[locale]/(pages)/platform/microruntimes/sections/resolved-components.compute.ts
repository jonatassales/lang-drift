import {
	hasPlayed,
	ticketPlan,
} from '@routes/platform/microruntimes/content/microruntimes-content'

export function resolvedComponents(playedCount: number) {
	return ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
}
