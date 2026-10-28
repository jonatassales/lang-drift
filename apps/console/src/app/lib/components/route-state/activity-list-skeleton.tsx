'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from '@ui/skeleton'

export function ActivityListSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="activity-list-skeleton"
		>
			<div className="grid h-[65vh] min-h-80 content-start gap-0 overflow-hidden rounded-xl border border-border">
				{['first', 'second', 'third', 'fourth'].map((key) => (
					<div key={key} className="grid gap-5 border-b border-border p-5">
						<Skeleton className="h-4 w-56 max-w-full" />
						<div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
							{['runtime', 'cause', 'state', 'time'].map((field) => (
								<Skeleton key={field} className="h-9 w-full" />
							))}
						</div>
					</div>
				))}
			</div>
		</section>
	)
}
