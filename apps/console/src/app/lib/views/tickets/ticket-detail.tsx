'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { StatusLabel } from '@components/status-label/status-label'
import { useResource } from '@data/resource.hook'
import type { TicketDetailSnapshot } from '@data/types'
import { Card } from '@ui/card'

import { PhaseWidget } from './phase-widget'
import { phaseReached, ticketPhases } from './ticket-phases'

interface TicketDetailProps {
	id: string
	initial: Promise<TicketDetailSnapshot>
}

// One ticket with every phase it has reached, in order, as the website's
// panels all open; a phase it has not reached is one muted line.
export function TicketDetail(props: TicketDetailProps) {
	const { id, initial } = props

	const t = useTranslations('tickets')
	const resource = useResource<TicketDetailSnapshot>(
		`tickets/${id}`,
		10000,
		initial,
	)
	const ticket = resource.data?.ticket

	if (!ticket) {
		return resource.error ? (
			<ErrorNotice error={resource.error} retry={resource.refresh} />
		) : null
	}

	return (
		<>
			<PageHeading title={`${ticket.number} · ${ticket.title}`}>
				<StatusLabel value={ticket.status} />
			</PageHeading>
			<div className="grid gap-4">
				{ticketPhases.map((phase) =>
					phaseReached(ticket, phase) ? (
						<Card className="p-5 sm:p-6" key={phase}>
							<PhaseWidget phase={phase} ticket={ticket} />
						</Card>
					) : (
						<p className="m-0 px-1 text-sm text-muted-foreground" key={phase}>
							{t('notReached', { phase: t(`labels.${phase}`) })}
						</p>
					),
				)}
			</div>
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.refresh} />
			) : null}
		</>
	)
}
