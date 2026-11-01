'use client'

import { Button } from '@langdrift/web/ui/button'
import { useTranslations } from 'next-intl'
import { Suspense, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { mutation } from '@data/api'
import { useResource } from '@data/resource.hook'
import type { LifecycleEvent, Runtime } from '@data/types'
import { RecommendationPanel } from '@views/runtimes/recommendation-panel'

import { EventHistory } from './event-history'

interface RuntimeHistoryProps {
	initial: Promise<{ runtime: Runtime }>
	id: string
	events: Promise<{ events: LifecycleEvent[] }>
}

export function RuntimeHistory(props: RuntimeHistoryProps) {
	const { id, initial, events } = props

	const t = useTranslations('ui')
	const resource = useResource<{ runtime: Runtime }>(
		`runtimes/${id}`,
		0,
		initial,
	)
	const [eventInitial, setEventInitial] = useState(events)
	const [revision, setRevision] = useState(0)
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	return (
		<>
			<PageHeading title={t('history')} description={id}>
				<Button
					disabled={pending}
					onClick={() => {
						if (!window.confirm(t('purgeConfirm'))) return
						setError(null)
						startTransition(async () => {
							try {
								await mutation(`runtimes/${id}/history`, 'DELETE')

								setEventInitial(
									Promise.resolve({
										events: [],
									}),
								)
								setRevision((current) => current + 1)
							} catch (failure) {
								setError(failure)
							}
						})
					}}
					type="button"
					variant="destructive"
				>
					{t('purge')}
				</Button>
			</PageHeading>
			<div className="grid gap-5">
				{error ? <ErrorNotice error={error} /> : null}
				<Suspense fallback={<ActivityListSkeleton />}>
					<EventHistory
						key={`events-${revision}`}
						runtimeId={id}
						initial={eventInitial}
					/>
				</Suspense>
				{resource.data ? (
					<RecommendationPanel
						productId={resource.data.runtime.productId}
						key={`recommendations-${revision}`}
						endpoint={`runtimes/${id}/recommendations`}
					/>
				) : null}
			</div>
		</>
	)
}
