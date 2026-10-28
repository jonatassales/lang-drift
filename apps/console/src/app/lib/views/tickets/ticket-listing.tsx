'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { ActivityRowSkeleton } from '@components/route-state/activity-row-skeleton'
import { usePagedResource } from '@data/paged-resource.hook'
import type { Page, Ticket, TicketPhase } from '@data/types'
import { Button } from '@ui/button'
import { VirtualList } from '@ui/virtual-list'

import { PhaseWidget } from './phase-widget'

interface TicketListingProps {
	initial: Promise<Page<Ticket>>
	/** The phase the listing shows; it lists every ticket that reached it. */
	phase: TicketPhase
	productId: string
}

const listLabels: Record<TicketPhase, string> = {
	ticket: 'tickets',
	analysis: 'analysis',
	plan: 'plans',
	workflow: 'workflows',
}

// One phase's listing: every ticket that has reached the phase, newest
// first, each as that phase's widget, loading older tickets as the list
// scrolls to its end.
export function TicketListing(props: TicketListingProps) {
	const { initial, phase, productId } = props

	const t = useTranslations('ui')
	const resource = usePagedResource<Ticket>(
		`products/${productId}/tickets?reached=${phase}`,
		initial,
	)

	return (
		<>
			<div className="mb-4 flex justify-end">
				<Button onClick={resource.refresh} disabled={resource.loading}>
					{t('refresh')}
				</Button>
			</div>
			{resource.items.length ? (
				<VirtualList
					items={resource.items}
					itemKey={(ticket) => ticket.id}
					renderItem={(ticket) => (
						<PhaseWidget listed phase={phase} ticket={ticket} />
					)}
					onEnd={resource.next}
					label={t(listLabels[phase])}
				/>
			) : null}
			{resource.loading ? <ActivityRowSkeleton /> : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.retry} />
			) : null}
			{!resource.loading && !resource.error && !resource.items.length ? (
				<p className="text-sm text-muted-foreground">{t('empty')}</p>
			) : null}
			{resource.hasMore ? (
				<div className="mt-4">
					<Button disabled={resource.loading} onClick={resource.next}>
						{t('loadMore')}
					</Button>
				</div>
			) : null}
		</>
	)
}
