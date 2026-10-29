import { notFound } from 'next/navigation'
import { Suspense } from 'react'

import { ticketTitle } from '@app/metadata/entity-title'
import { TicketDetailSkeleton } from '@components/route-state/ticket-detail-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type { TicketDetailSnapshot } from '@data/types'
import { TicketDetail } from '@views/tickets/ticket-detail'

interface PageProps {
	params: Promise<{
		id: string
		ticketId: string
	}>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { ticketId } = await params

	const title = await ticketTitle(ticketId)

	return { title }
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id, ticketId } = await params

	if (![id, ticketId].every(isEntityId)) notFound()

	return (
		<Suspense fallback={<TicketDetailSkeleton />}>
			<TicketDetail
				id={ticketId}
				initial={readConsole<TicketDetailSnapshot>(`tickets/${ticketId}`)}
			/>
		</Suspense>
	)
}
