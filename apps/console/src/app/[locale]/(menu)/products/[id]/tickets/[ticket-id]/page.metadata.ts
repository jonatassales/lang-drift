import { ticketTitle } from '@data/entity-title'

interface TicketRouteParams {
	id: string
	'ticket-id': string
}

interface PageProps {
	params: Promise<TicketRouteParams>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { 'ticket-id': ticketId } = await params

	const title = await ticketTitle(ticketId)

	return { title }
}
