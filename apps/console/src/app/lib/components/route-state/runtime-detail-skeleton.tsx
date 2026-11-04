'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from 'web/ui/skeleton'

import { HeadingSkeleton } from './heading-skeleton'

export function RuntimeDetailSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
			data-skeleton="runtime-detail-skeleton"
		>
			<HeadingSkeleton />
			<Skeleton className="h-5 w-44" />
			<Skeleton className="mt-4 h-[512px] w-full rounded-xl" />
			<div className="mt-6 grid gap-4 sm:grid-cols-3">
				{['cpu', 'ram', 'llm'].map((key) => (
					<Skeleton key={key} className="h-28 w-full" />
				))}
			</div>
		</section>
	)
}
