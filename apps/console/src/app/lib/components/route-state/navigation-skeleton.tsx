'use client'

import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

import { Navigation } from '@components/navigation/navigation'

export function NavigationSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="navigation-skeleton"
		>
			<Navigation products={[]} />
			<div className="mt-4 grid gap-3 border-t border-border px-3 pt-4">
				<Skeleton className="h-3 w-16" />
				<Skeleton className="h-12 w-full" />
				{['runtimes', 'executions', 'history', 'config'].map((key) => (
					<Skeleton key={key} className="h-5 w-28" />
				))}
			</div>
		</section>
	)
}
