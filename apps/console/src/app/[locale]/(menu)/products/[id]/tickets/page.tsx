import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { TicketListSkeleton } from '@components/route-state/ticket-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { TicketListing } from '@views/tickets/ticket-listing'

export { generateMetadata } from './page.metadata'

interface PageProps {
	params: Promise<{ id: string }>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading title={t('tickets')} description={t('ticketsIntro')} />
			<Suspense fallback={<TicketListSkeleton />}>
				<TicketListing
					initial={readConsole(`products/${id}/tickets?reached=ticket`)}
					phase="ticket"
					productId={id}
				/>
			</Suspense>
		</>
	)
}
