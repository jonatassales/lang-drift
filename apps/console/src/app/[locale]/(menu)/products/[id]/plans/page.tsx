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

	const [{ id }, translate] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading
				title={translate('tickets')}
				description={translate('ticketsIntro')}
			/>
			<Suspense fallback={<TicketListSkeleton />}>
				<TicketListing
					phase="plan"
					productId={id}
					initial={readConsole(`tickets?productId=${encodeURIComponent(id)}`)}
				/>
			</Suspense>
		</>
	)
}
