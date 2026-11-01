'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

import { TicketEntrySkeleton } from './ticket-entry-skeleton'

const entries = ['first', 'second', 'third', 'fourth', 'fifth', 'sixth']

export function TicketListSkeleton() {
	const t = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={t('loading')}
			data-skeleton="ticket-list-skeleton"
		>
			<div className="mb-4 flex justify-end">
				<Skeleton className="h-9 w-20" />
			</div>
			<div className="grid grid-cols-1 gap-3 lg:gap-2">
				{entries.map((key) => (
					<TicketEntrySkeleton key={key} />
				))}
			</div>
		</section>
	)
}
