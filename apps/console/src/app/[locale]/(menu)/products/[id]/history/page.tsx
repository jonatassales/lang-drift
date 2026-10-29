import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { projectTitle } from '@app/metadata/entity-title'
import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { HistoryOverview } from '@views/history/history-overview'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata(props: PageProps) {
	const { params } = props

	const { id } = await params

	const title = await projectTitle(id, 'history')

	return { title }
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading title={t('history')} description={t('activityIntro')} />
			<Suspense fallback={<ActivityListSkeleton />}>
				<HistoryOverview
					productId={id}
					initial={readConsole(`products/${id}/history`)}
				/>
			</Suspense>
		</>
	)
}
