import { useId } from 'react'

import { StatusLabel } from '@components/status-label/status-label'
import type { Ticket, TicketPhase } from '@data/types'
import { Link } from '@i18n/navigation'
import { CaretRight } from '@ui/icons'

import { PhaseSummary } from './phase-summary'
import { TicketCreated } from './ticket-created'
import { listedPhase, ticketPath } from './ticket-phases'

interface TicketRowProps {
	/** The listing the row sits in, which picks the summary it shows. */
	listing: TicketPhase
	ticket: Ticket
}

// A listing entry from lg: the ticket as one compact row, with its number,
// its title over its slug, its status, the listed phase's summary and the
// date it was created. Beside the sidebar the summary is a third line under
// the slug; from xl it moves under the status and the date takes a column
// of its own. The whole row opens the ticket.
export function TicketRow(props: TicketRowProps) {
	const { listing, ticket } = props

	const numberId = useId()
	const titleId = useId()
	const phase = listedPhase(ticket, listing)

	return (
		<article className="relative hidden grid-cols-[3rem_minmax(0,1fr)_9rem_1rem] items-center gap-x-4 gap-y-0.5 rounded-xl border border-border bg-card px-4 py-3 transition-colors hover:bg-muted lg:grid xl:grid-cols-[3rem_minmax(0,1fr)_minmax(0,22rem)_8.5rem_1rem]">
			<span
				className="col-start-1 row-span-3 row-start-1 font-mono text-[13px] tabular-nums xl:row-span-2"
				id={numberId}
			>
				{ticket.number}
			</span>
			<h2
				className="col-start-2 row-start-1 m-0 truncate text-[14px] font-medium"
				id={titleId}
			>
				{ticket.title}
			</h2>
			<p className="col-start-2 row-start-2 m-0 truncate font-mono text-[12px] text-muted-foreground">
				{ticket.slug}
			</p>
			{phase ? (
				<p className="col-start-2 row-start-3 m-0 truncate font-mono text-[12px] xl:col-start-3 xl:row-start-2">
					<PhaseSummary phase={phase} ticket={ticket} />
				</p>
			) : null}
			<span className="col-start-3 row-start-1 justify-self-end xl:justify-self-start">
				<StatusLabel value={ticket.status} />
			</span>
			<TicketCreated
				className="col-start-3 row-start-2 justify-self-end text-[12px] text-muted-foreground tabular-nums xl:col-start-4 xl:row-span-2 xl:row-start-1"
				createdAt={ticket.createdAt}
			/>
			<CaretRight
				aria-hidden
				className="col-start-4 row-span-3 row-start-1 text-muted-foreground xl:col-start-5 xl:row-span-2"
				size={14}
			/>
			<Link
				aria-labelledby={`${numberId} ${titleId}`}
				className="absolute inset-0 rounded-xl"
				href={ticketPath(ticket, listing)}
				prefetch={false}
			/>
		</article>
	)
}
