import type { MouseEvent } from 'react'

import {
	budgetCents,
	workflowDelivery,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import type {
	RunLine,
	RunMessageId,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.types'
import { runHasPlayed } from '@routes/platform/microruntimes/content/run-has-played.compute'
import type { WorkflowFigures } from '@routes/platform/microruntimes/content/workflow-figures.compute'
import { Eyebrow } from 'web/ui/eyebrow'
import { GitPullRequest } from 'web/ui/icons/server'
import { RevealValue } from 'web/ui/reveal-value'
import { cn } from 'web/utilities'

import { pullRequestState } from './pull-request-state.compute'
import type {
	TicketCardCopy,
	TicketCardFormat,
} from './ticket-card-format.compute'

const figure = 'font-mono text-[13px] tabular-nums'

type HasPlayed = (message: RunMessageId) => boolean

// What the run left as it plays: the reservation once the gate opens, the
// observed spend once the usage record closes and what the reservation
// released at the end, then the pull request, its checks and the Safe Seat
// state, which shows whether the floor was met and never a share. A run the
// gate paused shows the decision it needs and that nothing was paid.
// Every figure stays a dash until its line has printed. The simulated pull
// request reads as a link but opens nothing.
interface WorkflowResultProps {
	copy: TicketCardCopy
	figures: WorkflowFigures
	format: TicketCardFormat
	playedCount: number
	script: readonly RunLine[]
}

export function WorkflowResult(props: WorkflowResultProps) {
	const { copy, figures, format, playedCount, script } = props

	function handleAnchorClick(event: MouseEvent<HTMLAnchorElement>) {
		event.preventDefault()
	}

	const { workflow } = copy
	const played: HasPlayed = (message) =>
		runHasPlayed(script, playedCount, message)

	if (played('gatePaused')) {
		return (
			<div className="grid gap-1">
				<p className="m-0 text-[12px] text-status-warning">
					{workflow.needsDecision}
				</p>
				<p className="m-0 text-[11px] text-muted-foreground">
					{format.message('workflow.reserves', {
						budget: format.money(budgetCents),
						reserved: format.money(figures.reserved),
					})}{' '}
					· {workflow.nothingPaid}
				</p>
			</div>
		)
	}

	const cells = [
		{
			id: 'reserved',
			shown: played('gateOpened'),
			value: figures.reserved,
		},
		{
			id: 'observed',
			shown: played('completed'),
			value: figures.observed,
		},
		{
			id: 'released',
			shown: played('readyForReview') || played('readyToMerge'),
			value: figures.released,
		},
	] as const
	const safeSeatMet = played('safeSeatMet')

	return (
		<div className="grid gap-2.5">
			<dl className="m-0 grid grid-cols-3 gap-x-4">
				{cells.map((cell) => (
					<div className="grid gap-1" key={cell.id}>
						<Eyebrow as="dt">{workflow.result[cell.id]}</Eyebrow>
						<dd className="m-0">
							<RevealValue
								className={figure}
								labels={copy}
								state={cell.shown ? 'set' : 'unset'}
							>
								{format.money(cell.value)}
							</RevealValue>
						</dd>
					</div>
				))}
			</dl>
			<div className="flex flex-wrap items-baseline gap-x-3 gap-y-0.5 border-t border-border pt-2.5">
				<Eyebrow>{workflow.pullRequest}</Eyebrow>
				{played('pullRequest') ? (
					<>
						{/* biome-ignore lint/a11y/useValidAnchor: the simulated pull request reads as a link and opens nothing */}
						<a
							className="inline-flex items-center gap-1.5 font-mono text-[12px] text-foreground underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
							href="#"
							onClick={handleAnchorClick}
						>
							#{workflowDelivery.pullRequest} ·
							<GitPullRequest
								aria-hidden="true"
								className="shrink-0 text-status-success"
								size={14}
							/>
							{workflowDelivery.hash}
						</a>
						{played('checks') && (
							<span
								className={
									played('checksPassed')
										? 'text-[11px] text-status-success'
										: 'text-[11px] text-muted-foreground'
								}
							>
								{pullRequestState(workflow, played)}
							</span>
						)}
					</>
				) : (
					<RevealValue
						className="font-mono text-[12px]"
						labels={copy}
						state="unset"
					/>
				)}
			</div>
			<div className="flex flex-wrap items-baseline gap-x-3 gap-y-0.5 border-t border-border pt-2.5">
				<Eyebrow>{workflow.safeSeat}</Eyebrow>
				<RevealValue
					className={cn(
						'text-[11px]',
						safeSeatMet ? 'text-status-success' : 'text-muted-foreground',
					)}
					labels={copy}
					state={safeSeatMet ? 'set' : 'unset'}
				>
					{workflow.safeSeatMet}
				</RevealValue>
			</div>
		</div>
	)
}
