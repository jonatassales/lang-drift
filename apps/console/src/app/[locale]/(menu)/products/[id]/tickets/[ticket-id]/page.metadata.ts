import { ticketTitle } from '@data/entity-title'

interface PageProps {
	params: Promise<{ id: string; 'ticket-id': string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { 'ticket-id': ticketId } = await params

	const title = await ticketTitle(ticketId)

	return { title }
}
