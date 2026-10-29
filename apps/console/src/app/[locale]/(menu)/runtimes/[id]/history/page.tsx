import { notFound } from 'next/navigation'
import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import { RuntimeHistory } from '@views/history/runtime-history'

interface PageProps {
	params: Promise<{ id: string }>
}
export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('history') }
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { id } = await params

	if (!isEntityId(id)) notFound()

	return (
		<Suspense fallback={<ActivityListSkeleton />}>
			<RuntimeHistory
				id={id}
				initial={readConsole(`runtimes/${id}`)}
				events={readConsole(`runtimes/${id}/history`)}
			/>
		</Suspense>
	)
}
