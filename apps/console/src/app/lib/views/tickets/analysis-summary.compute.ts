import type { TicketAnalysis } from '@data/types'

import { keptEpisodes } from './kept-episodes.compute.ts'
import type { AnalysisSummary } from './ticket-phases.types.ts'

// The analysis's header line: the latest thing it has settled, from the
// documents read to the suggestion, or nothing yet.
export function analysisSummary(analysis: TicketAnalysis): AnalysisSummary {
	if (analysis.model && analysis.effort) {
		return {
			kind: 'suggestion',
			model: analysis.model.name,
			effort: analysis.effort,
		}
	}
	if (analysis.episodes && analysis.episodes.length > 0) {
		return {
			kind: 'episodes',
			...keptEpisodes(analysis.episodes),
		}
	}
	if (analysis.documentsRead.length > 0) {
		return {
			kind: 'documents',
			count: analysis.documentsRead.length,
		}
	}

	return { kind: 'unknown' }
}
