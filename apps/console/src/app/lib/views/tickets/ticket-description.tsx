import { useTranslations } from 'next-intl'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket } from '@data/types'

import { TicketCreated } from './ticket-created'
import { phaseProgress } from './ticket-phases'

interface TicketDescriptionProps {
	ticket: Ticket
}

// Where the ticket stands, as labels over values under the page heading:
// its status, the phase it is in out of the four, and when it was created.
export function TicketDescription(props: TicketDescriptionProps) {
	const { ticket } = props

	const translate = useTranslations('tickets')
	const progress = phaseProgress(ticket)
	const items = [
		{
			detail: <StatusLabel value={ticket.status} />,
			id: 'status',
			term: translate('ticket.status'),
		},
		{
			detail: translate('ticket.phaseProgress', {
				phase: translate(`labels.${ticket.phase}`),
				reached: progress.reached,
				total: progress.total,
			}),
			id: 'phase',
			term: translate('ticket.phase'),
		},
		{
			detail: <TicketCreated createdAt={ticket.createdAt} withYear />,
			id: 'created',
			term: translate('ticket.created'),
		},
	]

	return (
		<dl className="m-0 grid grid-cols-2 gap-x-8 gap-y-4 border-t border-border pt-5 sm:grid-cols-3">
			{items.map((item) => (
				<div className="grid content-start gap-1.5" key={item.id}>
					<dt className="text-[11px] text-muted-foreground">{item.term}</dt>
					<dd className="m-0 text-[14px] tabular-nums">{item.detail}</dd>
				</div>
			))}
		</dl>
	)
}
