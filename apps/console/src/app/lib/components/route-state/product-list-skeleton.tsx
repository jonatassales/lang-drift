'use client'

import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

import { HeadingSkeleton } from './heading-skeleton'

export function ProductListSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
			data-skeleton="product-list-skeleton"
		>
			<HeadingSkeleton />
			<div className="grid gap-4 md:grid-cols-2">
				{['first', 'second'].map((key) => (
					<div
						key={key}
						className="grid gap-4 rounded-xl border border-border p-5"
					>
						<Skeleton className="h-5 w-44" />
						<Skeleton className="h-4 w-full" />
						<Skeleton className="h-8 w-24" />
					</div>
				))}
			</div>
		</section>
	)
}
