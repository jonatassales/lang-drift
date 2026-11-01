import type { RuntimeDetailSnapshot, TicketDetailSnapshot } from '@data/types'
import { ticketSummary } from '@views/tickets/ticket-phases'

const ticketSections = ['tickets', 'analysis', 'plans', 'workflows']

export interface BreadcrumbItem {
	href: string
	key: string
	label: string
}

interface BreadcrumbOptions {
	base: string
	detail: RuntimeDetailSnapshot | null
	id: string
	section?: string
	segments: string[]
	ticket: TicketDetailSnapshot | null
	translate: (key: string) => string
}

export function breadcrumbItems(options: BreadcrumbOptions): BreadcrumbItem[] {
	const { base, detail, id, section, segments, ticket, translate } = options

	if (detail) {
		const items: BreadcrumbItem[] = [
			{
				href: base,
				key: 'runtimes',
				label: translate('runtimes'),
			},
			{
				href: `/runtimes/${id}`,
				key: 'runtime',
				label: detail.runtime.name,
			},
		]
		if (section === 'configuration' || section === 'history') {
			items.push({
				href: `/runtimes/${id}/${section}`,
				key: section,
				label: translate(section),
			})
		}

		return items
	}

	if (section === 'config') {
		return [
			{
				href: `${base}/config/secrets`,
				key: 'secrets',
				label: translate('secrets'),
			},
		]
	}

	const isNewRuntime = section === 'runtimes' && segments[3] === 'new'
	if (isNewRuntime) {
		return [
			{
				href: base,
				key: 'runtimes',
				label: translate('runtimes'),
			},
			{
				href: `${base}/runtimes/new`,
				key: 'new',
				label: translate('createRuntime'),
			},
		]
	}

	const isTicketSection = Boolean(section && ticketSections.includes(section))
	if (isTicketSection && section) {
		const items: BreadcrumbItem[] = [
			{
				href: `${base}/${section}`,
				key: section,
				label: translate(section),
			},
		]
		if (ticket) {
			items.push({
				href: `${base}/tickets/${ticket.ticket.id}`,
				key: 'ticket',
				label: ticketSummary(ticket.ticket),
			})
		}

		return items
	}

	const knownSection = section === 'executions' || section === 'history'
	const fallbackKey = knownSection && section ? section : 'runtimes'

	return [
		{
			href: section ? `${base}/${section}` : base,
			key: section || 'runtimes',
			label: translate(fallbackKey),
		},
	]
}
