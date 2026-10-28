import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { readConsole } from '@data/server'
import { TicketListing } from '@views/tickets/ticket-listing'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('analysis'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	return (
		<>
			<PageHeading title={t('analysis')} description={t('analysisIntro')} />
			<Suspense fallback={<ActivityListSkeleton />}>
				<TicketListing
					initial={readConsole(`products/${id}/tickets?reached=analysis`)}
					phase="analysis"
					productId={id}
				/>
			</Suspense>
		</>
	)
}
