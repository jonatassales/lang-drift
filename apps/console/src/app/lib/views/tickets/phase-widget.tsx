import type { Ticket, TicketPhase } from '@data/types'

import { AnalysisWidget } from './analysis-widget'
import { PlanWidget } from './plan-widget'
import { TicketWidget } from './ticket-widget'
import { WorkflowWidget } from './workflow-widget'

interface PhaseWidgetProps {
	phase: TicketPhase
	ticket: Ticket
}

// The widget for one phase of a ticket that has reached it.
export function PhaseWidget(props: PhaseWidgetProps) {
	const { phase, ticket } = props

	switch (phase) {
		case 'ticket':
			return <TicketWidget ticket={ticket} />
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
