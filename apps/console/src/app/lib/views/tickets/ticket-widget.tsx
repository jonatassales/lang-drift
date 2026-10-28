import { DetailList } from '@langdrift/web/ui/detail-list'
import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { useTranslations } from 'next-intl'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket } from '@data/types'
import { Link } from '@i18n/navigation'

import { ticketDocumentIds, ticketSummary } from './ticket-phases'
import { TicketWidgetFrame } from './ticket-widget-frame'
import { WidgetSection } from './widget-section'

interface TicketWidgetProps {
	/** Set in a listing, so the card names and links its ticket. */
	listed?: boolean
	ticket: Ticket
}

// The ticket as the website's Ticket panel draws it: its number, slug,
// status and title, then the artifacts it links.
export function TicketWidget(props: TicketWidgetProps) {
	const { listed = false, ticket } = props

	const t = useTranslations('tickets')

	return (
		<TicketWidgetFrame
			linkLine={false}
			phase="ticket"
			summary={ticketSummary(ticket)}
			ticket={listed ? ticket : undefined}
		>
			<div className="grid gap-2">
				<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
					<strong className="font-mono text-[40px] leading-none font-medium tracking-[-0.03em] max-sm:text-[32px]">
						{listed ? (
							<Link
								aria-label={t('openTicket', { number: ticket.number })}
								className="text-foreground underline-offset-4 hover:underline"
								href={`/products/${ticket.productId}/tickets/${ticket.id}`}
							>
								{ticket.number}
							</Link>
						) : (
							ticket.number
						)}
					</strong>
					<span className="font-mono text-[13px] text-muted-foreground">
						{ticket.slug}
					</span>
					<StatusLabel value={ticket.status} />
				</p>
				<p className="m-0 text-[15px] leading-snug font-medium">
					{ticket.title}
				</p>
				<WidgetSection>
					<Eyebrow>{t('ticket.links')}</Eyebrow>
					<DetailList
						className="mt-1"
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
				</WidgetSection>
			</div>
		</TicketWidgetFrame>
	)
}
