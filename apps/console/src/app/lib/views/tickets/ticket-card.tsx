import { useId } from 'react'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket, TicketPhase } from '@data/types'

import { PhaseSummary } from './phase-summary'
import { TicketCreated } from './ticket-created'
import { TicketEntryLink } from './ticket-entry-link'
import { listedPhase } from './ticket-phases'

interface TicketCardProps {
	/** The listing the card sits in, which picks the summary it shows. */
	listing: TicketPhase
	ticket: Ticket
}

// A listing entry below lg: the ticket as one card, with its number, the
// date it was created and its status, its title and slug, and a footer with
// the listed phase's summary. The whole card opens the ticket.
export function TicketCard(props: TicketCardProps) {
	const { listing, ticket } = props

	const numberId = useId()
	const titleId = useId()
	const phase = listedPhase(ticket, listing)

	return (
		<article className="relative grid grid-cols-1 gap-1.5 rounded-xl border border-border bg-card p-4 transition-colors hover:bg-muted lg:hidden">
			<div className="flex items-center justify-between gap-3">
				<span className="flex min-w-0 items-baseline gap-2">
					<span className="font-mono text-[13px] tabular-nums" id={numberId}>
						{ticket.number}
					</span>
					<TicketCreated
						className="truncate text-[12px] text-muted-foreground tabular-nums"
						createdAt={ticket.createdAt}
					/>
				</span>
				<StatusLabel value={ticket.status} />
			</div>
			<h2
				className="m-0 line-clamp-2 text-[15px] leading-snug font-medium"
				id={titleId}
			>
				{ticket.title}
			</h2>
			<p className="m-0 truncate font-mono text-[12px] text-muted-foreground">
				{ticket.slug}
			</p>
			{phase ? (
				<p className="m-0 mt-2 border-t border-border pt-3 font-mono text-[12px] leading-relaxed">
					<PhaseSummary phase={phase} ticket={ticket} />
				</p>
			) : null}
			<TicketEntryLink
				labelledBy={`${numberId} ${titleId}`}
				listing={listing}
				ticket={ticket}
			/>
		</article>
	)
}
