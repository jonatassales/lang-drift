import type { TicketWorkflow } from '@data/types'

import type { WorkflowSummary } from './ticket-phases.types.ts'

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
