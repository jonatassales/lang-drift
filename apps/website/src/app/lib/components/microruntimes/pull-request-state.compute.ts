import type { RunMessageId } from '@data/microruntimes-workflow'

import type { TicketCardCopy } from './ticket-card-format'

type HasPlayed = (message: RunMessageId) => boolean

export function pullRequestState(
	workflow: TicketCardCopy['workflow'],
	played: HasPlayed,
) {
	if (played('readyForReview')) return workflow.readyForReview
	if (played('readyToMerge')) return workflow.readyToMerge
	if (played('checksPassed')) return workflow.checksPassed

	return workflow.checksRunning
}
