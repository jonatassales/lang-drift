import {
	headroom,
	revealed,
	suggestedTotals,
	ticketPlan,
} from '@routes/platform/microruntimes/content/microruntimes-content'
import { Eyebrow } from 'web/ui/eyebrow'
import { ShareBar } from 'web/ui/share-bar'
import { cn } from 'web/utilities'

import { CardSection } from './card-section'
import { CardValue, type CardValueState, revealState } from './card-value'
import { ComponentsManifest } from './components-manifest'
import { resolvedComponents } from './resolved-components.compute'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

export { planSummary } from './plan-summary.compute'

interface PlanCardProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	playedCount: number
}

interface ReportCell {
	detail: string
	id: string
	state: CardValueState
	term: string
	tone?: 'warning'
}

const figure = 'font-mono tabular-nums'
const fade =
	'motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0'

// The plan as a receipt: the figures it proposes stepping up as the
// components resolve, the share of the envelope they take, the components
// themselves, then the report: what bounds the plan, what it leaves, what
// it would create, what runtime.md cannot bound, the price it cannot give
// and the validation it passed, closed by the evidence it leaves. Every
// figure is a proposal pending the team's authorization, which the card's
// only live region announces once, when the plan is drafted.
export function PlanCard(props: PlanCardProps) {
	const { copy, format, playedCount } = props

	const { plan } = copy
	const bounded = ticketPlan.bounded
	const resolved = resolvedComponents(playedCount)
	const totals = suggestedTotals(resolved)
	const totalsState = revealState(
		resolved.length > 0,
		revealed(playedCount, 'componentsSkeleton'),
	)
	const sized = revealed(playedCount, 'suggested')
	const drafted = revealed(playedCount, 'drafted')
	const envelopeRows = [
		{
			id: 'cpu',
			label: plan.envelope.cpu,
			bounded: format.cpu(bounded.cpuMillis),
			share: totals.cpuMillis / bounded.cpuMillis,
			suggested: format.cpu(totals.cpuMillis),
		},
		{
			id: 'memory',
			label: plan.envelope.memory,
			bounded: format.memory(bounded.memoryMiB),
			share: totals.memoryMiB / bounded.memoryMiB,
			suggested: format.memory(totals.memoryMiB),
		},
	]
	const reportCells: ReportCell[] = [
		{
			id: 'bounded',
			term: plan.report.bounded,
			state: revealState(revealed(playedCount, 'bounded')),
			detail: format.figures(bounded),
		},
		{
			id: 'headroom',
			term: plan.report.headroom,
			state: revealState(sized),
			detail: format.figures(headroom(bounded, totals)),
		},
		{
			id: 'resource',
			term: plan.report.suggestedResource,
			state: revealState(sized),
			detail: format.message('plan.deployments', {
				count: ticketPlan.components.length,
			}),
		},
		{
			id: 'limits',
			term: plan.report.limits,
			state: revealState(revealed(playedCount, 'limits')),
			detail: plan.notRuntimeFields,
		},
		{
			id: 'price',
			term: plan.report.price,
			state: revealState(revealed(playedCount, 'priceBasis')),
			detail: plan.priceBasis,
			tone: 'warning',
		},
		{
			id: 'validation',
			term: plan.report.validation,
			state: revealState(revealed(playedCount, 'validated')),
			detail: plan.validated,
		},
	]

	return (
		<div className="grid">
			<div className="grid gap-1.5">
				<Eyebrow>{plan.proposed}</Eyebrow>
				<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
					<CardValue
						className={`${figure} text-[28px] leading-none font-medium tracking-[-0.02em]`}
						labels={copy}
						state={totalsState}
						width="w-48"
					>
						{format.figures(totals)}
					</CardValue>
					{drafted && (
						<span className={`${fade} text-[12px] text-status-warning`}>
							{plan.pendingAuthorization}
						</span>
					)}
				</p>
			</div>

			<div className="mt-3 grid gap-2">
				{envelopeRows.map((row) => (
					<div className="grid gap-1" key={row.id}>
						<span className="flex items-baseline justify-between gap-3 text-[12px]">
							<span>{row.label}</span>
							<CardValue
								className={figure}
								labels={copy}
								state={totalsState}
								width="w-32"
							>
								{format.message('plan.ofBounded', {
									bounded: row.bounded,
									suggested: row.suggested,
								})}
							</CardValue>
						</span>
						<ShareBar share={totalsState === 'set' ? row.share : 0} />
					</div>
				))}
			</div>

			<ComponentsManifest
				copy={copy}
				format={format}
				playedCount={playedCount}
			/>
			{revealed(playedCount, 'crossChecked') && (
				<p className={`${fade} m-0 mt-1 text-[11px] text-muted-foreground`}>
					{format.message('plan.crossCheck', {
						ports: ticketPlan.components
							.map((component) => component.port)
							.join(' · '),
					})}
				</p>
			)}

			<CardSection className="mt-4 pt-3">
				<dl className="m-0 grid grid-cols-2 gap-x-6 gap-y-2.5 sm:grid-cols-3">
					{reportCells.map((cell) => (
						<div className="grid gap-1" key={cell.id}>
							<Eyebrow as="dt">{cell.term}</Eyebrow>
							<dd
								className={cn(
									'm-0 text-[12px] leading-snug',
									cell.state === 'set' && cell.tone === 'warning'
										? 'text-status-warning'
										: 'text-foreground',
								)}
							>
								<CardValue
									className={figure}
									labels={copy}
									state={cell.state}
									width="w-24"
								>
									{cell.detail}
								</CardValue>
							</dd>
						</div>
					))}
				</dl>
				<div
					aria-atomic="true"
					aria-live="polite"
					className="mt-3 border-t border-border"
				>
					{drafted && (
						<p
							className={`${fade} m-0 flex flex-wrap items-baseline gap-x-3 pt-2.5`}
						>
							<Eyebrow>{plan.report.evidence}</Eyebrow>
							<code className="font-mono text-[12px] break-all">
								{format.message('plan.planPath', {
									revision: ticketPlan.draftRevision,
								})}
							</code>
						</p>
					)}
				</div>
			</CardSection>
		</div>
	)
}
