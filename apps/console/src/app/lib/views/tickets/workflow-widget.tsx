import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { GitPullRequest } from '@langdrift/web/ui/icons/server'
import { ShareBar } from '@langdrift/web/ui/share-bar'
import { cn } from '@langdrift/web/utilities'

import type { Ticket, TicketWorkflow } from '@data/types'

import { PhaseSummary } from './phase-summary'
import { type TicketFormat, useTicketFormat } from './ticket-format.hook'
import { TicketSection } from './ticket-section'
import { UnknownValue } from './unknown-value'
import { WidgetSection } from './widget-section'
import { WorkflowStageList } from './workflow-stage-list'

interface WorkflowWidgetProps {
	ticket: Ticket
	workflow: TicketWorkflow
}

const figure = 'font-mono text-[13px] tabular-nums'

// Where the pull request stands: ready for the owner's review or to merge.
function pullRequestState(workflow: TicketWorkflow, format: TicketFormat) {
	const { t } = format
	if (workflow.outcome === 'readyToMerge') return t('workflow.readyToMerge')

	return t('workflow.readyForReview')
}

// The workflow as the website's Workflow panel draws it, read-only: the
// stages with their progress, the model and effort, the estimated spend
// and what it reserves of the budget, and the result with the reserved,
// observed and released amounts and the pull request. A run the gate
// paused shows the decision it needs and that nothing was paid.
export function WorkflowWidget(props: WorkflowWidgetProps) {
	const { ticket, workflow } = props

	const format = useTicketFormat()
	const { t } = format
	const overBudget = workflow.reserveCents > workflow.budgetCents
	const reserves = t('workflow.reserves', {
		budget: format.money(workflow.budgetCents),
		reserved: format.money(workflow.reserveCents),
	})
	const resultCells = [
		{
			id: 'reserved',
			value: workflow.reserveCents,
		},
		{
			id: 'observed',
			value: workflow.observedCents,
		},
		{
			id: 'released',
			value: workflow.releasedCents,
		},
	] as const

	return (
		<TicketSection
			id="workflow"
			summary={<PhaseSummary phase="workflow" ticket={ticket} />}
			title={t('labels.workflow')}
		>
			<WorkflowStageList stages={workflow.stages} />

			<WidgetSection className="grid grid-cols-2 gap-x-6 gap-y-2">
				<div className="grid gap-1.5">
					<Eyebrow>{t('workflow.model')}</Eyebrow>
					<span className="grid gap-0.5">
						<span className="text-[14px]">{workflow.model.name}</span>
						<span className="font-mono text-[11px] text-muted-foreground">
							{workflow.model.id}
						</span>
					</span>
				</div>
				<div className="grid gap-1.5">
					<Eyebrow>{t('workflow.effort')}</Eyebrow>
					<span className="font-mono text-[14px]">
						{t(`analysis.efforts.${workflow.effort}`)}
					</span>
				</div>
			</WidgetSection>

			<p className="m-0 mt-2 text-[11px] text-muted-foreground">
				{t('workflow.modelPrice', {
					date: format.priceMonth(workflow.priceDate),
					input: format.money(workflow.model.inputCents),
					output: format.money(workflow.model.outputCents),
				})}
			</p>

			<WidgetSection className="grid gap-1.5">
				<Eyebrow>{t('workflow.estimate')}</Eyebrow>
				<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
					<strong
						className={cn(
							'font-mono text-[28px] leading-none font-medium tracking-[-0.02em] tabular-nums',
							overBudget && 'text-status-warning',
						)}
					>
						{format.money(workflow.estimateCents)}
					</strong>
					<span className="text-[11px] text-muted-foreground">
						{t('workflow.tokens', {
							input: format.tokens(workflow.inputTokensK),
							output: format.tokens(workflow.outputTokensK),
						})}
					</span>
				</p>
				<span
					className={cn('text-[12px]', overBudget && 'text-status-warning')}
				>
					{reserves}
				</span>
				<ShareBar
					share={workflow.reserveCents / workflow.budgetCents}
					tone={overBudget ? 'warning' : 'default'}
				/>
				{overBudget && (
					<span className="text-[11px] text-status-warning">
						{t('workflow.overBudget')}
					</span>
				)}
			</WidgetSection>

			<WidgetSection>
				{workflow.outcome === 'paused' ? (
					<div className="grid gap-1">
						<p className="m-0 text-[12px] text-status-warning">
							{t('workflow.needsDecision')}
						</p>
						<p className="m-0 text-[11px] text-muted-foreground">
							{reserves} · {t('workflow.nothingPaid')}
						</p>
					</div>
				) : (
					<div className="grid gap-2.5">
						<dl className="m-0 grid grid-cols-3 gap-x-4">
							{resultCells.map((cell) => (
								<div className="grid gap-1" key={cell.id}>
									<Eyebrow as="dt">{t(`workflow.result.${cell.id}`)}</Eyebrow>
									<dd className="m-0">
										{cell.value === null ? (
											<UnknownValue className={figure} />
										) : (
											<span className={figure}>{format.money(cell.value)}</span>
										)}
									</dd>
								</div>
							))}
						</dl>
						<div className="flex flex-wrap items-baseline gap-x-3 gap-y-0.5 border-t border-border pt-2.5">
							<Eyebrow>{t('workflow.pullRequest')}</Eyebrow>
							{workflow.pullRequest ? (
								<>
									<span className="inline-flex items-center gap-1.5 font-mono text-[12px]">
										#{workflow.pullRequest.number} ·
										<GitPullRequest
											aria-hidden="true"
											className="shrink-0 text-status-success"
											size={14}
										/>
										{workflow.pullRequest.commit}
									</span>
									<span className="text-[11px] text-status-success">
										{t('workflow.checksPassed', {
											passed: workflow.pullRequest.checksPassed,
											total: workflow.pullRequest.checksTotal,
										})}{' '}
										· {pullRequestState(workflow, format)}
									</span>
								</>
							) : (
								<UnknownValue className="font-mono text-[12px]" />
							)}
						</div>
					</div>
				)}
			</WidgetSection>
		</TicketSection>
	)
}
