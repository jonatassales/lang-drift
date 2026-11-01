import type { Ticket, TicketPhase } from '@data/types'
import { Link } from '@i18n/navigation'

import { ticketPath } from './ticket-phases'

interface TicketEntryLinkProps {
	/** The ids of the elements that name the entry: its number and title. */
	labelledBy: string
	listing: TicketPhase
	ticket: Ticket
}

// The one link of a listing entry, laid over the whole widget so the entry
// opens its ticket wherever it is pressed and draws one focus ring.
export function TicketEntryLink(props: TicketEntryLinkProps) {
	const { labelledBy, listing, ticket } = props

	return (
		<Link
			aria-labelledby={labelledBy}
			className="absolute inset-0 rounded-xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
			href={ticketPath(ticket, listing)}
			prefetch={false}
		/>
	)
}
