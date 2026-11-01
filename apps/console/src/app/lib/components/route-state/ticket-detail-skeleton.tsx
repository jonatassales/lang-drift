'use client'
import { Skeleton } from '@langdrift/web/ui/skeleton'
import { useTranslations } from 'next-intl'

const cells = ['status', 'phase', 'created']
const sections = ['artifacts', 'analysis', 'plan']

export function TicketDetailSkeleton() {
	const translate = useTranslations('ui')

	return (
		<section
			aria-busy="true"
			aria-label={translate('loading')}
			data-skeleton="ticket-detail-skeleton"
		>
			<div className="mb-6 grid gap-2">
				<Skeleton className="h-3 w-48" />
				<Skeleton className="h-7 w-[28rem] max-w-full" />
			</div>
			<div className="grid grid-cols-2 gap-x-8 gap-y-4 border-t border-border pt-5 sm:grid-cols-3">
				{cells.map((cell) => (
					<div className="grid gap-1.5" key={cell}>
						<Skeleton className="h-3 w-16" />
						<Skeleton className="h-5 w-28" />
					</div>
				))}
			</div>
			<div className="mt-8 grid gap-10">
				{sections.map((section) => (
					<div className="grid gap-4 border-t border-border pt-6" key={section}>
						<Skeleton className="h-4 w-32" />
						<Skeleton className="h-24 w-full" />
					</div>
				))}
			</div>
		</section>
	)
}
