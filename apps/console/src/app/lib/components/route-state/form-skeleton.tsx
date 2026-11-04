'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from 'web/ui/skeleton'

import { HeadingSkeleton } from './heading-skeleton'

export function FormSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
			data-skeleton="form-skeleton"
		>
			<HeadingSkeleton />
			<div className="grid max-w-2xl gap-6 rounded-xl border border-border p-6">
				{['name', 'source', 'config'].map((key) => (
					<div key={key} className="grid gap-2">
						<Skeleton className="h-3 w-32" />
						<Skeleton className="h-10 w-full" />
					</div>
				))}
				<Skeleton className="h-9 w-32" />
			</div>
		</section>
	)
}
