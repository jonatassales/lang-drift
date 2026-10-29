import { useTranslations } from 'next-intl'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket } from '@data/types'

import { TicketCreated } from './ticket-created'
import { phaseStep } from './ticket-phases'

interface TicketDescriptionProps {
	ticket: Ticket
}

// Where the ticket stands, as labels over values under the page heading:
// its status, the phase it is in out of the four, and when it was created.
export function TicketDescription(props: TicketDescriptionProps) {
	const { ticket } = props

	const t = useTranslations('tickets')
	const step = phaseStep(ticket)
	const items = [
		{
			detail: <StatusLabel value={ticket.status} />,
			id: 'status',
			term: t('ticket.status'),
		},
		{
			detail: t('ticket.phaseProgress', {
				phase: t(`labels.${ticket.phase}`),
				reached: step.reached,
				total: step.total,
			}),
			id: 'phase',
			term: t('ticket.phase'),
		},
		{
			detail: <TicketCreated createdAt={ticket.createdAt} withYear />,
			id: 'created',
			term: t('ticket.created'),
		},
	]

	return (
		<dl className="m-0 grid grid-cols-2 gap-x-8 gap-y-4 border-t border-border pt-5 sm:grid-cols-3">
			{items.map((item) => (
				<div className="grid content-start gap-1.5" key={item.id}>
					<dt className="text-[11px] text-muted-foreground">{item.term}</dt>
					<dd className="m-0 text-sm tabular-nums">{item.detail}</dd>
				</div>
			))}
		</dl>
	)
}
