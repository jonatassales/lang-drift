'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

export function FieldSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="field-skeleton"
		>
			<div className="grid gap-3">
				<Skeleton className="h-4 w-64 max-w-full" />
				<Skeleton className="h-10 w-full" />
				<Skeleton className="h-9 w-28" />
			</div>
		</section>
	)
}
