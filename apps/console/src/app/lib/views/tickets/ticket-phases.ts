import type {
	Resources,
	Ticket,
	TicketAnalysis,
	TicketComponent,
	TicketDocumentId,
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

export const ticketDocumentIds: readonly TicketDocumentId[] = [
	'prd',
	'designDoc',
	'spec',
	'guardrail',
]

// The ticket's header line: its number and slug.
export function ticketSummary(ticket: Pick<Ticket, 'number' | 'slug'>) {
	return `${ticket.number} · ${ticket.slug}`
}

// True once the ticket's phase is this one or a later one.
export function phaseReached(ticket: Ticket, phase: TicketPhase) {
	return ticketPhases.indexOf(ticket.phase) >= ticketPhases.indexOf(phase)
}

// The phases with a one-line summary of their own.
export type SummaryPhase = Exclude<TicketPhase, 'ticket'>

export const summaryPhases: readonly SummaryPhase[] = [
	'analysis',
	'plan',
	'workflow',
]

// How far the ticket has come: the phases it has reached out of all four.
export function phaseStep(ticket: Pick<Ticket, 'phase'>) {
	return {
		reached: ticketPhases.indexOf(ticket.phase) + 1,
		total: ticketPhases.length,
	}
}

// Where a listing entry opens: the ticket's page, scrolled to the phase's
// section when the entry belongs to that phase's listing.
export function ticketPath(
	ticket: Pick<Ticket, 'id' | 'productId'>,
	listing: TicketPhase,
) {
	const page = `/products/${ticket.productId}/tickets/${ticket.id}`
	if (listing === 'ticket') return page

	return `${page}#${listing}`
}

// The phase whose summary a listing entry shows: the listing's own phase,
// or, in the Tickets listing, the phase the ticket is in, which has none
// while the ticket is still in its Ticket phase.
export function listedPhase(
	ticket: Pick<Ticket, 'phase'>,
	listing: TicketPhase,
): SummaryPhase | null {
	const phase = listing === 'ticket' ? ticket.phase : listing
	if (phase === 'ticket') return null

	return phase
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

export interface PlanFigures {
	suggested: Resources | null
	headroom: Resources | null
	cpuShare: number
	memoryShare: number
}

// A component's share of the plan's CPU envelope, from 0 to 1.
export function componentShare(
	component: Pick<TicketComponent, 'cpuMillis'>,
	bounded: Resources,
) {
	return component.cpuMillis / bounded.cpuMillis
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
	| { kind: 'estimate' }

// The workflow's header line: how the run ended, or, until it has ended
// with a pull request and an observed spend, its estimate.
export function workflowSummary(workflow: TicketWorkflow): WorkflowSummary {
	if (workflow.outcome === 'paused') return { kind: 'paused' }
	if (workflow.outcome === 'running') return { kind: 'estimate' }
	if (!workflow.pullRequest || workflow.observedCents === null) {
		return { kind: 'estimate' }
	}

	return {
		kind: 'done',
		number: workflow.pullRequest.number,
		observedCents: workflow.observedCents,
		state: workflow.outcome,
	}
}
