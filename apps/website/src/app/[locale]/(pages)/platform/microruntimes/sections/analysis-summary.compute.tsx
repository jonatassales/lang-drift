import type { ReactNode } from 'react'

import {
	documentIds,
	episodes,
	revealed,
	ticketPlan,
} from '@routes/platform/microruntimes/content/microruntimes-content'
import { workflowModel } from '@routes/platform/microruntimes/content/microruntimes-workflow'

import { CardValue } from './card-value'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

const suggestedModel = workflowModel(ticketPlan.model)
const keptCount = episodes.filter(
	(episode) => episode.verdict === 'kept',
).length

export function analysisSummary(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	playedCount: number,
): ReactNode {
	if (revealed(playedCount, 'suggestion')) {
		return format.message('analysis.suggestionSummary', {
			effort: copy.analysis.efforts[ticketPlan.effort],
			model: suggestedModel.name,
		})
	}
	if (revealed(playedCount, 'verdicts')) {
		return format.message('analysis.episodesSummary', {
			kept: keptCount,
			total: episodes.length,
		})
	}
	if (revealed(playedCount, 'documentsRead')) {
		return format.message('analysis.documentsSummary', {
			count: documentIds.length,
		})
	}

	return <CardValue labels={copy} state="unset" />
}
