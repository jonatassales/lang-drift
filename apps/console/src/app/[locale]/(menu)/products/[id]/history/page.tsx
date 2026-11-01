import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityOverviewSkeleton } from '@components/route-state/activity-overview-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { HistoryOverview } from '@views/history/history-overview'

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
			<PageHeading title={t('history')} description={t('activityIntro')} />
			<Suspense fallback={<ActivityOverviewSkeleton />}>
				<HistoryOverview
					productId={id}
					initial={readConsole(`products/${id}/history`)}
				/>
			</Suspense>
		</>
	)
}
