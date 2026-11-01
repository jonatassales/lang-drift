'use client'

import { Button } from '@langdrift/web/ui/button'
import { ListEnd } from '@langdrift/web/ui/list-end'
import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { ActivityEntrySkeleton } from '@components/route-state/activity-entry-skeleton'
import { usePagedResource } from '@data/paged-resource.hook'
import type { Page, ProductEvent } from '@data/types'

import { HistoryEntry } from './history-entry'

interface HistoryOverviewProps {
	productId: string
	initial: Promise<Page<ProductEvent>>
}

// Every lifecycle event in the project, newest first, in the page's flow,
// each in an item of its own, loading older ones as the list's end comes
// into view.
export function HistoryOverview(props: HistoryOverviewProps) {
	const { productId, initial } = props

	const t = useTranslations('ui')
	const resource = usePagedResource<ProductEvent>(
		`products/${productId}/history`,
		initial,
	)

	return (
		<>
			<div className="mb-4 flex justify-end">
				<Button
					onClick={resource.refresh}
					disabled={resource.loading}
					type="button"
					variant="outline"
				>
					{t('refresh')}
				</Button>
			</div>
			{resource.items.length ? (
				<ol
					aria-label={t('history')}
					className="m-0 grid list-none grid-cols-1 gap-3 p-0"
				>
					{resource.items.map((item) => (
						<li key={item.event.id}>
							<HistoryEntry item={item} />
						</li>
					))}
				</ol>
			) : null}
			{resource.loading ? (
				<div
					aria-busy="true"
					aria-label={t('loading')}
					className="mt-3"
					role="status"
				>
					<ActivityEntrySkeleton />
				</div>
			) : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.retry} />
			) : null}
			{!resource.loading && !resource.error && !resource.items.length ? (
				<p className="text-sm text-muted-foreground">{t('empty')}</p>
			) : null}
			<ListEnd
				failed={Boolean(resource.error)}
				hasMore={resource.hasMore}
				label={t('loadMore')}
				loading={resource.loading}
				onEnd={resource.next}
			/>
		</>
	)
}
