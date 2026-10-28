import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { useTranslations } from 'next-intl'
import type { ReactNode } from 'react'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket, TicketPhase } from '@data/types'
import { Link } from '@i18n/navigation'

import { ticketSummary } from './ticket-phases'

interface TicketWidgetFrameProps {
	children: ReactNode
	phase: TicketPhase
	/** The panel's one line: its latest figure. */
	summary: ReactNode
	/** Set in a listing: the ticket the card belongs to, linked to its page. */
	ticket?: Ticket
	/** In a listing, the Ticket widget links its own number instead. */
	linkLine?: boolean
}

// A ticket widget as the website draws its panel: the phase's eyebrow and
// one-line summary over the full body, never collapsed. In a listing, a
// line above names the ticket, links to its page and gives its status.
export function TicketWidgetFrame(props: TicketWidgetFrameProps) {
	const { children, linkLine = true, phase, summary, ticket } = props

	const t = useTranslations('tickets')

	return (
		<article className="grid min-w-0 gap-3 text-left">
			{ticket && linkLine && (
				<div className="flex flex-wrap items-center justify-between gap-2">
					<Link
						className="font-mono text-[12px] text-foreground underline-offset-4 hover:underline"
						href={`/products/${ticket.productId}/tickets/${ticket.id}`}
					>
						{ticketSummary(ticket)}
					</Link>
					<StatusLabel value={ticket.status} />
				</div>
			)}
			<header className="flex min-w-0 items-center justify-between gap-4">
				<Eyebrow className="shrink-0">{t(`labels.${phase}`)}</Eyebrow>
				<span className="min-w-0 truncate font-mono text-[12px] text-muted-foreground">
					{summary}
				</span>
			</header>
			<div>{children}</div>
		</article>
	)
}
