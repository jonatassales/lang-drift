import { notFound } from 'next/navigation'

import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type { TicketDetailSnapshot } from '@data/types'
import { TicketDetail } from '@views/tickets/ticket-detail'

export { generateMetadata } from './page.metadata'

interface TicketRouteParams {
	id: string
	'ticket-id': string
}

interface PageProps {
	params: Promise<TicketRouteParams>
}

export default async function ProductTicketPage(props: PageProps) {
	const { params } = props

	const { id, 'ticket-id': ticketId } = await params

	if (!isEntityId(id) || !isEntityId(ticketId)) notFound()

	return (
		<TicketDetail
			id={ticketId}
			initial={readConsole<TicketDetailSnapshot>(`tickets/${ticketId}`)}
		/>
	)
}
