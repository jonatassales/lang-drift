'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { ActivityRowSkeleton } from '@components/route-state/activity-row-skeleton'
import { usePagedResource } from '@data/paged-resource.hook'
import type { Page, ProductExecution } from '@data/types'
import { Button } from '@ui/button'
import { VirtualList } from '@ui/virtual-list'

import { ExecutionEntry } from './execution-entry'

interface ExecutionOverviewProps {
	productId: string
	initial: Promise<Page<ProductExecution>>
}

export function ExecutionOverview(props: ExecutionOverviewProps) {
	const { productId, initial } = props

	const t = useTranslations('ui')
	const resource = usePagedResource<ProductExecution>(
		`products/${productId}/executions`,
		initial,
	)

	return (
		<>
			<div className="mb-4 flex justify-end">
				<Button onClick={resource.refresh} disabled={resource.loading}>
					{t('refresh')}
				</Button>
			</div>
			{resource.items.length ? (
				<VirtualList
					items={resource.items}
					itemKey={(item) => item.execution.id}
					renderItem={(item) => <ExecutionEntry item={item} />}
					onEnd={resource.next}
					label={t('executions')}
				/>
			) : null}
			{resource.loading ? <ActivityRowSkeleton /> : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.retry} />
			) : null}
			{!resource.loading && !resource.error && !resource.items.length ? (
				<p className="text-sm text-muted-foreground">{t('empty')}</p>
			) : null}
			{resource.hasMore ? (
				<div className="mt-4">
					<Button disabled={resource.loading} onClick={resource.next}>
						{t('loadMore')}
					</Button>
				</div>
			) : null}
		</>
	)
}
