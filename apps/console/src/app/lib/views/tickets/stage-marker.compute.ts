import type { TicketStage } from '@data/types'

// The mark a workflow stage carries: a person writes it by hand, a person
// approves the agent's work, or the agent runs it alone. Merge carries its
// own lock instead.
export function stageMarker(stage: TicketStage): 'person' | 'approval' | null {
	if (stage.drivenBy === 'person') return 'person'
	if (stage.human && stage.id !== 'merge') return 'approval'

	return null
}
