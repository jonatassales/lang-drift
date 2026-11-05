'use client'

import { useTranslations } from 'next-intl'
import { Suspense, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { ActivityListSkeleton } from '@components/route-state/activity-list-skeleton'
import { mutation } from '@data/api'
import type {
	LifecycleEventsEnvelope,
	RuntimeEnvelope,
} from '@data/read-envelopes.types'
import { useResource } from '@data/resource.hook'
import { RecommendationPanel } from '@views/runtimes/recommendation-panel'
import { Button } from 'web/ui/button'

import { EventHistory } from './event-history'

interface RuntimeHistoryProps {
	initial: Promise<RuntimeEnvelope>
	id: string
	events: Promise<LifecycleEventsEnvelope>
}

export function RuntimeHistory(props: RuntimeHistoryProps) {
	const { id, initial, events } = props

	const [eventInitial, setEventInitial] = useState(events)
	const [revision, setRevision] = useState(0)
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	const resource = useResource<RuntimeEnvelope>(`runtimes/${id}`, 0, initial)
	const translate = useTranslations('ui')

	function handlePurge() {
		if (!window.confirm(translate('purgeConfirm'))) return
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
	}

	return (
		<>
			<PageHeading title={translate('history')} description={id}>
				<Button
					disabled={pending}
					onClick={handlePurge}
					type="button"
					variant="destructive"
				>
					{translate('purge')}
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
