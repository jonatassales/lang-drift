import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityOverviewSkeleton } from '@components/route-state/activity-overview-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { ExecutionOverview } from '@views/executions/execution-overview'

export { generateMetadata } from './page.metadata'

interface ProductRouteParams {
	id: string
}

interface PageProps {
	params: Promise<ProductRouteParams>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const [{ id }, translate] = await Promise.all([params, getTranslations('ui')])

	if (!isEntityId(id)) notFound()

	return (
		<>
			<PageHeading
				title={translate('executions')}
				description={translate('activityIntro')}
			/>
			<Suspense fallback={<ActivityOverviewSkeleton />}>
				<ExecutionOverview
					productId={id}
					initial={readConsole(`products/${id}/executions`)}
				/>
			</Suspense>
		</>
	)
}
