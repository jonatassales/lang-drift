import { DetailList } from '@langdrift/web/ui/detail-list'

import {
	hasPlayed,
	revealed,
	suggestedTotals,
	ticketPlan,
} from '@data/microruntimes-content'

import { CardSection } from './card-section'
import { CardValue, revealState } from './card-value'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

interface PlanReportProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

const figure = 'font-mono tabular-nums'

// The closing report of the plan: the bounded figures, the suggested figures
// stepping with the components, the suggested materialization, and a final
// row pending the team's authorization. The final row and the draft line
// sit in the card's only live region, so they are announced once, when the
// plan is drafted.
export function PlanReport(props: PlanReportProps) {
	const { copy, format, playedCount } = props

	const resolved = ticketPlan.components.filter((component) =>
		hasPlayed(playedCount, component.id),
	)
	const suggested = suggestedTotals(ticketPlan.components)
	const drafted = revealed(playedCount, 'report')
	const rows = [
		{
			id: 'bounded',
			term: copy.report.bounded,
			state: revealState(revealed(playedCount, 'bounded')),
			detail: format.figures(ticketPlan.bounded),
		},
		{
			id: 'suggested',
			term: copy.report.suggested,
			state: revealState(resolved.length > 0),
			detail: format.figures(suggestedTotals(resolved)),
		},
		{
			id: 'resource',
			term: copy.report.suggestedResource,
			state: revealState(drafted),
			detail: format.message('deployments', {
				count: ticketPlan.components.length,
			}),
		},
	]

	return (
		<CardSection>
			<DetailList
				items={rows.map((row) => ({
					detail: (
						<CardValue
							className={figure}
							labels={copy}
							state={row.state}
							width="w-40"
						>
							{row.detail}
						</CardValue>
					),
					id: row.id,
					term: row.term,
				}))}
			/>
			<div
				aria-atomic="true"
				aria-live="polite"
				className="grid gap-2.5 border-t border-border"
			>
				{drafted && (
					<>
						<DetailList
							items={[
								{
									detail: (
										<span className={`${figure} text-[14px] font-medium`}>
											{format.message('finalDetail', {
												cpu: format.formatNumber.format(suggested.cpuMillis),
												memory: format.formatNumber.format(suggested.memoryMiB),
												pending: copy.pendingAuthorization,
											})}
										</span>
									),
									id: 'final',
									term: copy.report.final,
								},
							]}
						/>
						<p className="m-0 font-mono text-[11px] text-muted-foreground">
							{format.message('planPath', {
								revision: ticketPlan.draftRevision,
							})}
						</p>
					</>
				)}
			</div>
		</CardSection>
	)
}
