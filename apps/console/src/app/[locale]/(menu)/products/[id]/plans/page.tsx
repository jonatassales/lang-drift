import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { TicketListing } from '@views/tickets/ticket-listing'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('plans'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading title={t('plans')} description={t('plansIntro')} />
			<Suspense fallback={<ActivityListSkeleton />}>
				<TicketListing
					initial={readConsole(`products/${id}/tickets?reached=plan`)}
					phase="plan"
					productId={id}
				/>
			</Suspense>
		</>
	)
}
