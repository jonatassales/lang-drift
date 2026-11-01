import type { Ticket, TicketPhase } from '@data/types'

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
