import { cn } from '@langdrift/web/utilities'
import { type ChangeEvent, useId } from 'react'

import {
	type HumanStageId,
	isRequiredHumanStage,
	type RunLine,
	stageProgress,
	stageSkills,
	type WorkflowChoices,
	workflowStageIds,
} from '@data/microruntimes-workflow'

import { ReadMark } from './read-mark'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'

interface WorkflowStagesProps {
	copy: TicketCardCopy
	format: TicketCardFormat
	human: WorkflowChoices['human']
	locked: boolean
	onToggle: (stage: HumanStageId) => void
	/** The run's lines and how many have played, while a run is shown. */
	run: { playedCount: number; script: readonly RunLine[] } | null
}

const skillText = 'font-mono text-[10px] text-muted-foreground'

// The delivery stages the workflow runs, each with the skills it uses and a
// Human in the loop checkbox. Planning always has a person, so its box stays
// checked and unavailable, with a note on hover or focus saying why; the dev
// ticket is driven by the person, not approved, and merge is always the
// owner's, so neither has a box. While the run plays, the stage at work
// takes the primary tone and a finished stage its read mark; a reviewed run
// ends waiting for the owner.
export function WorkflowStages(props: WorkflowStagesProps) {
	const { copy, format, human, locked, onToggle, run } = props

	const noteId = useId()

	function handleToggleStage(event: ChangeEvent<HTMLInputElement>) {
		const stage = event.currentTarget.dataset.stage as HumanStageId | undefined
		if (stage) {
			onToggle(stage)
		}
	}

	const { workflow } = copy

	return (
		<fieldset className="m-0 min-w-0 border-0 p-0">
			<legend className="mb-1.5 p-0 font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase">
				{workflow.humanInTheLoop}
			</legend>
			<ul className="m-0 grid list-none gap-x-4 gap-y-1 p-0 sm:grid-cols-2">
				{workflowStageIds.map((stage) => {
					const skills = stageSkills[stage]

					if (stage === 'merge') {
						return (
							<li className="flex min-h-7 items-center gap-2" key={stage}>
								<ReadMark label={workflow.mergeLocked} />
								<span className="text-[12px]">{workflow.stages.merge}</span>
								<span className={skillText} aria-hidden="true">
									{workflow.mergeLocked}
								</span>
							</li>
						)
					}

					const progress = run
						? stageProgress(run.script, run.playedCount, stage)
						: 'pending'

					if (stage === 'dev-ticket') {
						return (
							<li
								className={cn(
									'flex min-h-7 items-center gap-2 text-[12px]',
									progress === 'active' && 'text-primary',
								)}
								key={stage}
							>
								<span>{workflow.stages['dev-ticket']}</span>
								<span className={skillText}>{workflow.personDriven}</span>
								{progress === 'done' && (
									<ReadMark label={workflow.stages['dev-ticket']} />
								)}
								{progress === 'active' && (
									<span className="text-[10px] text-status-warning">
										{workflow.personWriting}
									</span>
								)}
							</li>
						)
					}

					const required = isRequiredHumanStage(stage)
					const unavailable = locked || required
					const requiredNoteId = `${noteId}-${stage}`
					// A person waits on a stage while it runs, except review: the
					// agent's code review runs first, and the owner's review begins
					// once the pull request is ready.
					const ownerTurn = stage === 'review' ? 'done' : 'active'
					const waiting = human[stage] && progress === ownerTurn

					return (
						<li className="group relative" key={stage}>
							<label
								className={cn(
									'flex min-h-7 cursor-pointer items-center gap-2 text-[12px]',
									progress === 'active' && 'text-primary',
									unavailable && 'cursor-default',
								)}
							>
								<input
									aria-describedby={required ? requiredNoteId : undefined}
									aria-disabled={unavailable || undefined}
									aria-label={format.message('workflow.humanStage', {
										stage: workflow.stages[stage],
									})}
									checked={human[stage]}
									className={cn(
										'size-3.5 shrink-0 cursor-[inherit] accent-primary focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
										required && 'cursor-not-allowed opacity-50',
									)}
									data-stage={stage}
									onChange={handleToggleStage}
									type="checkbox"
								/>
								<span>{workflow.stages[stage]}</span>
								{skills && <span className={skillText}>{skills}</span>}
								{progress === 'done' && !waiting && (
									<ReadMark label={workflow.stages[stage]} />
								)}
								{waiting && (
									<span className="text-[10px] text-status-warning">
										{workflow.waiting}
									</span>
								)}
							</label>
							{required && (
								<span
									className="pointer-events-none invisible absolute top-full left-0 z-10 mt-0.5 w-max max-w-72 rounded-md border border-border bg-foreground px-2 py-1 text-[11px] leading-snug text-background opacity-0 motion-safe:transition-opacity group-focus-within:visible group-focus-within:opacity-100 group-hover:visible group-hover:opacity-100"
									id={requiredNoteId}
									role="tooltip"
								>
									{format.message('workflow.requiredStage', {
										stage: workflow.stages[stage],
									})}
								</span>
							)}
						</li>
					)
				})}
			</ul>
		</fieldset>
	)
}
