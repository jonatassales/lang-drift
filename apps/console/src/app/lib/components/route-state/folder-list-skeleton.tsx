'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

export function FolderListSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
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
