import type {
	Resources,
	Ticket,
	TicketAnalysis,
	TicketEpisode,
	TicketPhase,
	TicketPlan,
	TicketWorkflow,
} from '@data/types'

export const ticketPhases: readonly TicketPhase[] = [
	'ticket',
	'analysis',
	'plan',
	'workflow',
]

// True once the ticket's phase is this one or a later one.
export function phaseReached(ticket: Ticket, phase: TicketPhase) {
	return ticketPhases.indexOf(ticket.phase) >= ticketPhases.indexOf(phase)
}

export type AnalysisSummary =
	| { kind: 'suggestion'; model: string; effort: string }
	| { kind: 'episodes'; kept: number; total: number }
	| { kind: 'documents'; count: number }
	| { kind: 'unknown' }

export function keptEpisodes(episodes: readonly TicketEpisode[]) {
	return {
		kept: episodes.filter((episode) => episode.verdict === 'kept').length,
		total: episodes.length,
	}
}

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
	if (analysis.episodes.length > 0) {
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

export interface PlanFigures {
	suggested: Resources | null
	headroom: Resources | null
	cpuShare: number
	memoryShare: number
}

// The plan's totals, what they leave in the envelope and the share they
// take of it. Without suggested totals there is no headroom, never a zero
// one, and the bars stay empty.
export function planFigures(plan: TicketPlan): PlanFigures {
	const { bounded, suggested } = plan
	if (!suggested) {
		return {
			suggested: null,
			headroom: null,
			cpuShare: 0,
			memoryShare: 0,
		}
	}

	return {
		suggested,
		headroom: {
			cpuMillis: bounded.cpuMillis - suggested.cpuMillis,
			memoryMiB: bounded.memoryMiB - suggested.memoryMiB,
		},
		cpuShare: suggested.cpuMillis / bounded.cpuMillis,
		memoryShare: suggested.memoryMiB / bounded.memoryMiB,
	}
}

export type WorkflowSummary =
	| {
			kind: 'done'
			number: number
			observedCents: number
			state: 'readyForReview' | 'readyToMerge'
	  }
	| { kind: 'paused' }
	| { kind: 'running' }

// The workflow's header line: how the run ended, or that it is running.
export function workflowSummary(workflow: TicketWorkflow): WorkflowSummary {
	if (workflow.outcome === 'paused') return { kind: 'paused' }
	if (workflow.outcome === 'running') return { kind: 'running' }
	if (!workflow.pullRequest || workflow.observedCents === null) {
		return { kind: 'running' }
	}

	return {
		kind: 'done',
		number: workflow.pullRequest.number,
		observedCents: workflow.observedCents,
		state: workflow.outcome,
	}
}
