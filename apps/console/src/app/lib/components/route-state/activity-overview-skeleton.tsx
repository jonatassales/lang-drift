'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from '@ui/skeleton'

import { ActivityEntrySkeleton } from './activity-entry-skeleton'

const entries = ['first', 'second', 'third', 'fourth', 'fifth', 'sixth']

export function ActivityOverviewSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="activity-overview-skeleton"
		>
			<div className="mb-4 flex justify-end">
				<Skeleton className="h-9 w-20" />
			</div>
			<div className="grid grid-cols-1 gap-3">
				{entries.map((key) => (
					<ActivityEntrySkeleton key={key} />
				))}
			</div>
		</section>
	)
}
