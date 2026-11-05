import type { RunMessageId } from '@routes/platform/microruntimes/content/microruntimes-workflow.data'

import type { TicketCardCopy } from './ticket-card-format.factory'

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
