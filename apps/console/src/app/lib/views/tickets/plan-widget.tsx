import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { ShareBar } from '@langdrift/web/ui/share-bar'
import { cn } from '@langdrift/web/utilities'

import type { Ticket, TicketPlan } from '@data/types'

import { ComponentRow } from './component-row'
import { componentsLine } from './components-line.compute'
import { PhaseSummary } from './phase-summary'
import { useTicketFormat } from './ticket-format.hook'
import { componentShare, planFigures } from './ticket-phases'
import { TicketSection } from './ticket-section'
import { UnknownValue } from './unknown-value'
import { WidgetSection } from './widget-section'

interface PlanWidgetProps {
	plan: TicketPlan
	ticket: Ticket
}

interface ReportCell {
	detail: string | null
	id: string
	term: string
	tone?: 'warning'
}

const figure = 'font-mono tabular-nums'

// The plan as the website's Plan panel draws it: the figures it proposes,
// the share of the envelope they take, the components with their marks and
// shares, then the report: what bounds the plan, what it leaves, what it
// would create, what runtime.md cannot bound, the price it cannot give and
// the validation, closed by the evidence it leaves.
export function PlanWidget(props: PlanWidgetProps) {
	const { plan, ticket } = props

	const format = useTicketFormat()
	const { t } = format
	const figures = planFigures(plan)
	const drafted = plan.suggested !== null
	const bounded = plan.bounded
	const envelopeRows = [
		{
			id: 'cpu',
			label: t('plan.envelope.cpu'),
			bounded: format.cpu(bounded.cpuMillis),
			share: figures.cpuShare,
			suggested: figures.suggested && format.cpu(figures.suggested.cpuMillis),
		},
		{
			id: 'memory',
			label: t('plan.envelope.memory'),
			bounded: format.memory(bounded.memoryMiB),
			share: figures.memoryShare,
			suggested:
				figures.suggested && format.memory(figures.suggested.memoryMiB),
		},
	]
	const reportCells: ReportCell[] = [
		{
			id: 'bounded',
			term: t('plan.report.bounded'),
			detail: format.figures(bounded),
		},
		{
			id: 'headroom',
			term: t('plan.report.headroom'),
			detail: figures.headroom && format.figures(figures.headroom),
		},
		{
			id: 'resource',
			term: t('plan.report.suggestedResource'),
			detail: figures.suggested
				? t('plan.deployments', { count: plan.components.length })
				: null,
		},
		{
			id: 'limits',
			term: t('plan.report.limits'),
			detail: t('plan.notRuntimeFields'),
		},
		{
			id: 'price',
			term: t('plan.report.price'),
			detail: t('plan.priceBasis'),
			tone: 'warning',
		},
		{
			id: 'validation',
			term: t('plan.report.validation'),
			detail:
				plan.validationErrors === null
					? null
					: t('plan.validated', { count: plan.validationErrors }),
		},
	]

	return (
		<TicketSection
			id="plan"
			summary={<PhaseSummary phase="plan" ticket={ticket} />}
			title={t('labels.plan')}
		>
			<div className="grid gap-1.5">
				<Eyebrow>{t('plan.proposed')}</Eyebrow>
				<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
					{figures.suggested ? (
						<span
							className={`${figure} text-[28px] leading-none font-medium tracking-[-0.02em]`}
						>
							{format.figures(figures.suggested)}
						</span>
					) : (
						<UnknownValue
							className={`${figure} text-[28px] leading-none font-medium`}
						/>
					)}
					{drafted && (
						<span
							className={cn(
								'text-[12px]',
								plan.authorization === 'authorized'
									? 'text-status-success'
									: 'text-status-warning',
							)}
						>
							{format.authorization(plan)}
						</span>
					)}
				</p>
			</div>

			<div className="mt-3 grid gap-2">
				{envelopeRows.map((row) => (
					<div className="grid gap-1" key={row.id}>
						<span className="flex items-baseline justify-between gap-3 text-[12px]">
							<span>{row.label}</span>
							{row.suggested ? (
								<span className={figure}>
									{t('plan.ofBounded', {
										bounded: row.bounded,
										suggested: row.suggested,
									})}
								</span>
							) : (
								<UnknownValue className={figure} />
							)}
						</span>
						<ShareBar share={row.share} />
					</div>
				))}
			</div>

			<WidgetSection>
				<div className="grid grid-cols-[minmax(0,1fr)_auto_auto] gap-x-3">
					<Eyebrow>{t('plan.components')}</Eyebrow>
					<Eyebrow>{t('plan.columns.cpu')}</Eyebrow>
					<Eyebrow>{t('plan.columns.memory')}</Eyebrow>
				</div>
				<ul className="m-0 mt-1 list-none p-0">
					{plan.components.map((component) => (
						<ComponentRow
							component={component}
							cpu={format.cpu(component.cpuMillis)}
							key={component.path}
							kind={t(`plan.kinds.${component.kind}`)}
							memory={format.memory(component.memoryMiB)}
							share={componentShare(component, bounded)}
						/>
					))}
				</ul>
				<p className="m-0 mt-1 border-t border-border pt-2 text-[12px] text-muted-foreground">
					{componentsLine(
						plan,
						ticket.analysis?.specComponents ?? null,
						format,
					)}
				</p>
				{drafted && (
					<p className="m-0 mt-1 text-[11px] text-muted-foreground">
						{t('plan.crossCheck', {
							ports: plan.components
								.map((component) => component.port)
								.join(' · '),
						})}
					</p>
				)}
			</WidgetSection>

			<WidgetSection>
				<dl className="m-0 grid grid-cols-2 gap-x-6 gap-y-2.5 sm:grid-cols-3">
					{reportCells.map((cell) => (
						<div className="grid gap-1" key={cell.id}>
							<Eyebrow as="dt">{cell.term}</Eyebrow>
							<dd
								className={cn(
									'm-0 text-[12px] leading-snug',
									cell.tone === 'warning'
										? 'text-status-warning'
										: 'text-foreground',
								)}
							>
								{cell.detail === null ? (
									<UnknownValue className={figure} />
								) : (
									<span className={figure}>{cell.detail}</span>
								)}
							</dd>
						</div>
					))}
				</dl>
				{drafted && (
					<p className="m-0 mt-3 flex flex-wrap items-baseline gap-x-3 border-t border-border pt-2.5">
						<Eyebrow>{t('plan.report.evidence')}</Eyebrow>
						<code className="font-mono text-[12px] break-all">
							{t('plan.planPath', { revision: plan.revision })}
						</code>
					</p>
				)}
			</WidgetSection>
		</TicketSection>
	)
}
