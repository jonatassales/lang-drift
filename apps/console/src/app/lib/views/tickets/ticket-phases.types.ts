import type { Resources, TicketDocumentId, TicketPhase } from '@data/types'

export const ticketPhases: readonly TicketPhase[] = [
	'ticket',
	'analysis',
	'plan',
	'workflow',
]

export const ticketDocumentIds: readonly TicketDocumentId[] = [
	'prd',
	'designDoc',
	'spec',
	'guardrail',
]

export type SummaryPhase = Exclude<TicketPhase, 'ticket'>

export const summaryPhases: readonly SummaryPhase[] = [
	'analysis',
	'plan',
	'workflow',
]

export interface AnalysisSuggestionSummary {
	effort: string
	kind: 'suggestion'
	model: string
}

export interface AnalysisEpisodesSummary {
	kept: number
	kind: 'episodes'
	total: number
}

export interface AnalysisDocumentsSummary {
	count: number
	kind: 'documents'
}

export interface AnalysisUnknownSummary {
	kind: 'unknown'
}

export type AnalysisSummary =
	| AnalysisSuggestionSummary
	| AnalysisEpisodesSummary
	| AnalysisDocumentsSummary
	| AnalysisUnknownSummary

export interface PlanFigures {
	suggested: Resources | null
	headroom: Resources | null
	cpuShare: number
	memoryShare: number
}

export interface WorkflowDoneSummary {
	kind: 'done'
	number: number
	observedCents: number
	state: 'readyForReview' | 'readyToMerge'
}

export interface WorkflowPausedSummary {
	kind: 'paused'
}

export interface WorkflowEstimateSummary {
	kind: 'estimate'
}

export type WorkflowSummary =
	| WorkflowDoneSummary
	| WorkflowPausedSummary
	| WorkflowEstimateSummary
