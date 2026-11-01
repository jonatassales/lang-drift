import type { TicketWorkflow } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'

// Where the pull request stands: ready for the owner's review or to merge.
export function pullRequestState(
	workflow: TicketWorkflow,
	format: TicketFormat,
): string {
	const { translate } = format
	if (workflow.outcome === 'readyToMerge')
		return translate('workflow.readyToMerge')

	return translate('workflow.readyForReview')
}
