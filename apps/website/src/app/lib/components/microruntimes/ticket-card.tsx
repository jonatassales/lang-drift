import { documentIds, revealed, ticketPlan } from '@data/microruntimes-content'
import { DetailList } from 'web/ui/detail-list'
import { Eyebrow } from 'web/ui/eyebrow'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import type { TicketCardCopy } from './ticket-card-format'

export { ticketSummary } from './ticket-summary.compute'

interface TicketCardProps {
	copy: TicketCardCopy
	playedCount: number
}

const fileName = 'font-mono text-[12px]'

// The ticket as the agent reads it: its number, slug, status and title,
// then the artifacts it links, which land once the agent has read the
// ticket. The reads themselves belong to the analysis card.
export function TicketCard(props: TicketCardProps) {
	const { copy, playedCount } = props

	const linksState = revealState(revealed(playedCount, 'links'))

	return (
		<div className="grid gap-2">
			<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<strong className="font-mono text-[44px] leading-none font-medium tracking-[-0.03em] max-sm:text-[36px]">
					{ticketPlan.ticket.number}
				</strong>
				<span className="font-mono text-[13px] text-muted-foreground">
					{ticketPlan.ticket.slug}
				</span>
				<span className="rounded-full border border-border px-2 py-[3px] font-mono text-[11px]">
					{copy.ticket.status}
				</span>
			</p>
			<p className="m-0 text-[15px] leading-snug font-medium">
				{copy.ticket.title}
			</p>

			<CardSection>
				<Eyebrow>{copy.ticket.links}</Eyebrow>
				<DetailList
					className="mt-1"
					items={documentIds.map((id) => ({
						detail: (
							<CardValue
								className={fileName}
								labels={copy}
								state={linksState}
								width="w-56"
							>
								{ticketPlan.artifacts[id]}
							</CardValue>
						),
						id,
						term: copy.ticket.artifacts[id],
					}))}
				/>
			</CardSection>
		</div>
	)
}
