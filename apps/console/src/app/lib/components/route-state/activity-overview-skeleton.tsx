'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

import { ActivityEntrySkeleton } from './activity-entry-skeleton'

const entries = ['first', 'second', 'third', 'fourth', 'fifth', 'sixth']

export function ActivityOverviewSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
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
