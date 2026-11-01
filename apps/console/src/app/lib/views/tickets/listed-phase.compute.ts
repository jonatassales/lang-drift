import type { Ticket, TicketPhase } from '@data/types'

import type { SummaryPhase } from './ticket-phases.types.ts'

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
