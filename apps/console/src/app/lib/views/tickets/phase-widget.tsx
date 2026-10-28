import type { Ticket, TicketPhase } from '@data/types'

import { AnalysisWidget } from './analysis-widget'
import { TicketWidget } from './ticket-widget'

interface PhaseWidgetProps {
	/** Set in a listing, so the card names and links its ticket. */
	listed?: boolean
	phase: TicketPhase
	ticket: Ticket
}

// The widget for one phase of a ticket that has reached it.
export function PhaseWidget(props: PhaseWidgetProps) {
	const { listed = false, phase, ticket } = props

	switch (phase) {
		case 'ticket':
			return <TicketWidget listed={listed} ticket={ticket} />
		case 'analysis':
			return ticket.analysis ? (
				<AnalysisWidget
					analysis={ticket.analysis}
					listed={listed}
					ticket={ticket}
				/>
			) : null
		default:
			return null
	}
}
