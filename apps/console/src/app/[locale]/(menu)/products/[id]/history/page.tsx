import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityOverviewSkeleton } from '@components/route-state/activity-overview-skeleton'
import { readConsole } from '@data/server'
import { HistoryOverview } from '@views/history/history-overview'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('history'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

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
