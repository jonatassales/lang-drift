import type { Ticket } from '@data/types'

// The ticket's header line: its number and slug.
export function ticketSummary(ticket: Pick<Ticket, 'number' | 'slug'>) {
	return `${ticket.number} · ${ticket.slug}`
}
