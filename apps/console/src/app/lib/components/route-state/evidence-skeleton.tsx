'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

export function EvidenceSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="evidence-skeleton"
		>
			<div className="grid gap-3 rounded-lg border border-border p-4">
				<Skeleton className="h-4 w-48 max-w-full" />
				<Skeleton className="h-3 w-full" />
				<Skeleton className="h-3 w-3/4" />
			</div>
		</section>
	)
}
