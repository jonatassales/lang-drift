'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { useResource } from '@data/resource.hook'
import type { TicketDetailSnapshot } from '@data/types'

import { useFragmentLanding } from './fragment-landing.hook'
import { PhaseWidget } from './phase-widget'
import { TicketArtifacts } from './ticket-artifacts'
import { TicketDescription } from './ticket-description'
import { phaseReached, summaryPhases, ticketSummary } from './ticket-phases'

interface TicketDetailProps {
	id: string
	initial: Promise<TicketDetailSnapshot>
}

// One ticket as one document: its title under its number and slug, where
// it stands, the artifacts it links, then each phase it has reached as a
// section with the website panel's body; a phase it has not reached yet is
// one muted line.
export function TicketDetail(props: TicketDetailProps) {
	const { id, initial } = props

	const translate = useTranslations('tickets')
	const resource = useResource<TicketDetailSnapshot>(
		`tickets/${id}`,
		10000,
		initial,
	)
	const ticket = resource.data?.ticket
	useFragmentLanding(Boolean(ticket))

	if (!ticket) {
		return resource.error ? (
			<ErrorNotice error={resource.error} retry={resource.refresh} />
		) : null
	}

	return (
		<>
			<header className="mb-6 grid gap-2">
				<p className="m-0 font-mono text-[12px] text-muted-foreground tabular-nums">
					{ticketSummary(ticket)}
				</p>
				<h1 className="m-0 text-[24px] font-medium tracking-tight text-balance">
					{ticket.title}
				</h1>
			</header>
			<TicketDescription ticket={ticket} />
			<div className="mt-8 grid gap-10">
				<TicketArtifacts ticket={ticket} />
				{summaryPhases.map((phase) =>
					phaseReached(ticket, phase) ? (
						<PhaseWidget key={phase} phase={phase} ticket={ticket} />
					) : (
						<p
							className="m-0 border-t border-border pt-6 text-sm text-muted-foreground"
							key={phase}
						>
							{translate('notReached', { phase: translate(`labels.${phase}`) })}
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
