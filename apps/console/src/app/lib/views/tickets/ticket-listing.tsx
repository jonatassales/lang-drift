'use client'

import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { ActivityRowSkeleton } from '@components/route-state/activity-row-skeleton'
import { usePagedResource } from '@data/paged-resource.hook'
import type { Page, Ticket, TicketPhase } from '@data/types'
import { Button } from '@ui/button'
import { ListEnd } from '@ui/list-end'

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
// first, in the page's flow, loading older tickets as the list's end comes
// into view.
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
				<ol
					aria-label={t(listLabels[phase])}
					className="m-0 grid list-none gap-3 p-0"
				>
					{resource.items.map((ticket) => (
						<li
							className="rounded-xl border border-border bg-card p-4"
							key={ticket.id}
						>
							<PhaseWidget listed phase={phase} ticket={ticket} />
						</li>
					))}
				</ol>
			) : null}
			{resource.loading ? <ActivityRowSkeleton /> : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.retry} />
			) : null}
			{!resource.loading && !resource.error && !resource.items.length ? (
				<p className="text-sm text-muted-foreground">{t('empty')}</p>
			) : null}
			<ListEnd
				failed={Boolean(resource.error)}
				hasMore={resource.hasMore}
				label={t('loadMore')}
				loading={resource.loading}
				onEnd={resource.next}
			/>
		</>
	)
}
