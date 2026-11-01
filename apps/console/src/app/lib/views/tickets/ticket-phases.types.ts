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

export type AnalysisSummary =
	| { kind: 'suggestion'; model: string; effort: string }
	| { kind: 'episodes'; kept: number; total: number }
	| { kind: 'documents'; count: number }
	| { kind: 'unknown' }

export interface PlanFigures {
	suggested: Resources | null
	headroom: Resources | null
	cpuShare: number
	memoryShare: number
}

export type WorkflowSummary =
	| {
			kind: 'done'
			number: number
			observedCents: number
			state: 'readyForReview' | 'readyToMerge'
	  }
	| { kind: 'paused' }
	| { kind: 'estimate' }
