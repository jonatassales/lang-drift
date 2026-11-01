import type { Ticket } from '@data/types'

import { ticketPhases } from './ticket-phases.types.ts'

// How far the ticket has come: the phases it has reached out of all four.
export function phaseProgress(ticket: Pick<Ticket, 'phase'>) {
	return {
		reached: ticketPhases.indexOf(ticket.phase) + 1,
		total: ticketPhases.length,
	}
}
