import type { Ticket } from '@data/types'

import { AnalysisWidget } from './analysis-widget'
import { PlanWidget } from './plan-widget'
import type { SummaryPhase } from './ticket-phases'
import { WorkflowWidget } from './workflow-widget'

interface PhaseWidgetProps {
	phase: SummaryPhase
	ticket: Ticket
}

// The section for one phase of a ticket that has reached it.
export function PhaseWidget(props: PhaseWidgetProps) {
	const { phase, ticket } = props

	switch (phase) {
		case 'analysis':
			return ticket.analysis ? (
				<AnalysisWidget analysis={ticket.analysis} ticket={ticket} />
			) : null
		case 'plan':
			return ticket.plan ? (
				<PlanWidget plan={ticket.plan} ticket={ticket} />
			) : null
		case 'workflow':
			return ticket.workflow ? (
				<WorkflowWidget ticket={ticket} workflow={ticket.workflow} />
			) : null
	}
}
