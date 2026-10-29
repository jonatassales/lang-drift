import { DetailList } from '@langdrift/web/ui/detail-list'
import { useTranslations } from 'next-intl'

import type { Ticket } from '@data/types'

import { ticketDocumentIds } from './ticket-phases'
import { TicketSection } from './ticket-section'

interface TicketArtifactsProps {
	ticket: Ticket
}

// The artifacts the ticket links, by kind, as the website's Ticket panel
// lists them; the page header already carries its number, title and slug.
export function TicketArtifacts(props: TicketArtifactsProps) {
	const { ticket } = props

	const t = useTranslations('tickets')

	return (
		<TicketSection id="artifacts" title={t('ticket.links')}>
			<DetailList
				className="max-w-3xl"
				items={ticketDocumentIds.map((id) => ({
					detail: (
						<span className="font-mono text-[12px] break-all">
							{ticket.artifacts[id]}
						</span>
					),
					id,
					term: t(`ticket.artifacts.${id}`),
				}))}
			/>
		</TicketSection>
	)
}
