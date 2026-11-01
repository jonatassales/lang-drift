'use client'

import { Button } from '@langdrift/web/ui/button'
import { ListEnd } from '@langdrift/web/ui/list-end'
import { useTranslations } from 'next-intl'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { TicketEntrySkeleton } from '@components/route-state/ticket-entry-skeleton'
import { usePagedResource } from '@data/paged-resource.hook'
import type { Page, Ticket, TicketPhase } from '@data/types'

import { TicketCard } from './ticket-card'
import { TicketRow } from './ticket-row'

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
// first, in the page's flow, each as a card below lg and a row from lg,
// loading older tickets as the list's end comes into view.
export function TicketListing(props: TicketListingProps) {
	const { initial, phase, productId } = props

	const translate = useTranslations('ui')
	const resource = usePagedResource<Ticket>(
		`products/${productId}/tickets?reached=${phase}`,
		initial,
	)

	return (
		<>
			<div className="mb-4 flex justify-end">
				<Button
					onClick={resource.refresh}
					disabled={resource.loading}
					type="button"
					variant="outline"
				>
					{translate('refresh')}
				</Button>
			</div>
			{resource.items.length ? (
				<ol
					aria-label={translate(listLabels[phase])}
					className="m-0 grid list-none grid-cols-1 gap-3 p-0 lg:gap-2"
				>
					{resource.items.map((ticket) => (
						<li key={ticket.id}>
							<TicketCard listing={phase} ticket={ticket} />
							<TicketRow listing={phase} ticket={ticket} />
						</li>
					))}
				</ol>
			) : null}
			{resource.loading ? (
				<div
					aria-busy="true"
					aria-label={translate('loading')}
					className="mt-3 lg:mt-2"
					role="status"
				>
					<TicketEntrySkeleton />
				</div>
			) : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.retry} />
			) : null}
			{!resource.loading && !resource.error && !resource.items.length ? (
				<p className="text-sm text-muted-foreground">{translate('empty')}</p>
			) : null}
			<ListEnd
				failed={Boolean(resource.error)}
				hasMore={resource.hasMore}
				label={translate('loadMore')}
				loading={resource.loading}
				onEnd={resource.next}
			/>
		</>
	)
}
