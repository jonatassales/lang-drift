import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { RuntimeDetailSkeleton } from '@components/route-state/runtime-detail-skeleton'
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

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('tickets'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id, ticketId } = await params

	if (![id, ticketId].every(isEntityId)) notFound()

	return (
		<Suspense fallback={<RuntimeDetailSkeleton />}>
			<TicketDetail
				id={ticketId}
				initial={readConsole<TicketDetailSnapshot>(`tickets/${ticketId}`)}
			/>
		</Suspense>
	)
}
