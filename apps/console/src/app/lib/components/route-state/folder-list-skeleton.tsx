'use client'
import { useTranslations } from 'next-intl'

import { Skeleton } from '@ui/skeleton'

export function FolderListSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="folder-list-skeleton"
		>
			<div className="grid gap-2 py-3">
				{['first', 'second', 'third', 'fourth'].map((key) => (
					<div
						key={key}
						className="flex gap-3 rounded-lg border border-border p-3"
					>
						<Skeleton className="h-5 w-5" />
						<Skeleton className="h-5 w-44 max-w-full" />
					</div>
				))}
			</div>
		</section>
	)
}
