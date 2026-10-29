import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { ExecutionOverview } from '@views/executions/execution-overview'

interface PageProps {
	params: Promise<{ id: string }>
}

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return {
		title: t('executions'),
	}
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, t] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading title={t('executions')} description={t('activityIntro')} />
			<Suspense fallback={<ActivityListSkeleton />}>
				<ExecutionOverview
					productId={id}
					initial={readConsole(`products/${id}/executions`)}
				/>
			</Suspense>
		</>
	)
}
