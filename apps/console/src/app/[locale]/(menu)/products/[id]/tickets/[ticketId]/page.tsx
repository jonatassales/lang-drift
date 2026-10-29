import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { TicketDetailSkeleton } from '@components/route-state/ticket-detail-skeleton'
import { readConsole } from '@data/server'
import type { TicketDetailSnapshot } from '@data/types'
import { TicketDetail } from '@views/tickets/ticket-detail'

interface PageProps {
	params: Promise<{ ticketId: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('tickets'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { ticketId } = await params

	return (
		<Suspense fallback={<TicketDetailSkeleton />}>
			<TicketDetail
				id={ticketId}
				initial={readConsole<TicketDetailSnapshot>(
					`tickets/${encodeURIComponent(ticketId)}`,
				)}
			/>
		</Suspense>
	)
}
