'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from '@ui/skeleton'

import { HeadingSkeleton } from './heading-skeleton'

export function RuntimeListSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="runtime-list-skeleton"
		>
			<HeadingSkeleton />
			<div className="grid gap-3">
				{['first', 'second', 'third'].map((key) => (
					<div
						key={key}
						className="grid gap-4 rounded-xl border border-border bg-card p-4"
					>
						<Skeleton className="h-5 w-60 max-w-full" />
						<Skeleton className="h-3 w-44" />
						<div className="grid gap-3 sm:grid-cols-3">
							{['cpu', 'ram', 'llm'].map((field) => (
								<Skeleton key={field} className="h-8 w-full" />
							))}
						</div>
						<Skeleton className="h-3 w-24" />
					</div>
				))}
			</div>
		</section>
	)
}
