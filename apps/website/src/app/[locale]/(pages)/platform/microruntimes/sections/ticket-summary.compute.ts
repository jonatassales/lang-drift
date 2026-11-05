import { ticketPlan } from '@routes/platform/microruntimes/content/microruntimes-content.data'

export function ticketSummary() {
	return `${ticketPlan.ticket.number} · ${ticketPlan.ticket.slug}`
}
