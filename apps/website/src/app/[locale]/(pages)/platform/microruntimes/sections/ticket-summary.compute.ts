import { ticketPlan } from '@routes/platform/microruntimes/content/microruntimes-content'

export function ticketSummary() {
	return `${ticketPlan.ticket.number} · ${ticketPlan.ticket.slug}`
}
