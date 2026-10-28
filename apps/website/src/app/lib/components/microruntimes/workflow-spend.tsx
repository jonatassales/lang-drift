import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'

import type { PlanEffort } from '@data/microruntimes-content'
import {
	budgetCents,
	expectedTokens,
	type WorkflowFigures,
} from '@data/microruntimes-workflow'

import { ShareBar } from './share-bar'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

interface WorkflowSpendProps {
	copy: TicketCardCopy
	effort: PlanEffort
	figures: WorkflowFigures
	format: TicketCardFormat
}

// What the chosen run is expected to cost and what the gate would reserve
// of the budget the design doc set, with the tokens behind the estimate.
// A reservation over the budget takes the warning role: the gate would
// pause the run before any paid call.
export function WorkflowSpend(props: WorkflowSpendProps) {
	const { copy, effort, figures, format } = props

	const { workflow } = copy
	const tokens = expectedTokens[effort]
	const tone = figures.overBudget ? 'warning' : 'default'

	return (
		<div className="grid gap-1.5">
			<Eyebrow>{workflow.estimate}</Eyebrow>
			<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<strong
					className={cn(
						'font-mono text-[28px] leading-none font-medium tracking-[-0.02em] tabular-nums',
						figures.overBudget && 'text-status-warning',
					)}
				>
					{format.money(figures.estimate)}
				</strong>
				<span className="text-[11px] text-muted-foreground">
					{format.message('workflow.tokens', {
						input: format.tokens(tokens.inputThousands),
						output: format.tokens(tokens.outputThousands),
					})}
				</span>
			</p>
			<span
				className={cn(
					'text-[12px]',
					figures.overBudget && 'text-status-warning',
				)}
			>
				{format.message('workflow.reserves', {
					budget: format.money(budgetCents),
					reserved: format.money(figures.reserved),
				})}
			</span>
			<ShareBar share={figures.reserved / budgetCents} tone={tone} />
			{figures.overBudget && (
				<span className="text-[11px] text-status-warning">
					{workflow.overBudget}
				</span>
			)}
		</div>
	)
}
