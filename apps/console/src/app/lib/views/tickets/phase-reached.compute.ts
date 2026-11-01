import type { Ticket, TicketPhase } from '@data/types'

import { ticketPhases } from './ticket-phases.types.ts'

// True once the ticket's phase is this one or a later one.
export function phaseReached(ticket: Ticket, phase: TicketPhase) {
	return ticketPhases.indexOf(ticket.phase) >= ticketPhases.indexOf(phase)
}
