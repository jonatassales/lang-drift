import { ticketPlan } from '@data/microruntimes-content'

export function ticketSummary() {
	return `${ticketPlan.ticket.number} · ${ticketPlan.ticket.slug}`
}
