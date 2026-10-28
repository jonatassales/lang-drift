'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from '@ui/skeleton'

export function ActivityRowSkeleton() {
	const t = useTranslations('ui')

	return (
		<div
			role="status"
			className="mt-3 grid gap-4 rounded-xl border border-border p-5"
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="activity-row-skeleton"
		>
			<Skeleton className="h-4 w-56 max-w-full" />
			<div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
				{['runtime', 'cause', 'state', 'time'].map((field) => (
					<Skeleton key={field} className="h-9 w-full" />
				))}
			</div>
		</div>
	)
}
