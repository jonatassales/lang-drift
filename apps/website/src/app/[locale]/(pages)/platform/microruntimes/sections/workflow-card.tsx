import { useAtomValue, useSetAtom } from 'jotai'

import {
	workflowEfforts,
	workflowModels,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import { workflowFigures } from '@routes/platform/microruntimes/content/workflow-figures.compute'
import { workflowModel } from '@routes/platform/microruntimes/content/workflow-model.compute'
import {
	canPlayWorkflowAtom,
	chooseEffortAtom,
	chooseModelAtom,
	microruntimeSessionAtom,
	playWorkflowAtom,
	runInProgressAtom,
	runPlayedCountAtom,
	runScriptAtom,
	toggleHumanStageAtom,
	workflowChoicesAtom,
} from '@state'
import { CardSection } from 'web/ui/card-section'

import { playLabel } from './play-label.compute'
import type {
	TicketCardCopy,
	TicketCardFormat,
} from './ticket-card-format.compute'
import { WorkflowChoice } from './workflow-choice'
import { WorkflowResult } from './workflow-result'
import { WorkflowSpend } from './workflow-spend'
import { WorkflowStages } from './workflow-stages'

const playButton =
	'inline-flex min-h-[38px] cursor-pointer items-center gap-[7px] rounded-full bg-primary px-[17px] whitespace-nowrap text-primary-foreground focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary motion-safe:transition-[transform,filter] motion-safe:duration-200 hover:saturate-[1.18] motion-safe:hover:-translate-y-px aria-disabled:cursor-not-allowed aria-disabled:opacity-60 aria-disabled:hover:translate-y-0'

// The session's last card: the run the plan leads to, summarized as money.
// The visitor picks the human stages, a Claude model and an effort, sees the
// expected spend against the budget the team set, and plays the workflow.
// The card then follows the run the terminal plays, from the gate to the
// pull request.
interface WorkflowCardProps {
	copy: TicketCardCopy
	format: TicketCardFormat
}

export function WorkflowCard(props: WorkflowCardProps) {
	const { copy, format } = props

	const session = useAtomValue(microruntimeSessionAtom)
	const choices = useAtomValue(workflowChoicesAtom)
	const script = useAtomValue(runScriptAtom)
	const runPlayedCount = useAtomValue(runPlayedCountAtom)
	const running = useAtomValue(runInProgressAtom)
	const canPlay = useAtomValue(canPlayWorkflowAtom)
	const chooseModel = useSetAtom(chooseModelAtom)
	const chooseEffort = useSetAtom(chooseEffortAtom)
	const toggleHumanStage = useSetAtom(toggleHumanStageAtom)
	const play = useSetAtom(playWorkflowAtom)

	function handlePlay() {
		play()
	}

	const { workflow } = copy
	const model = workflowModel(choices.model)
	const showsRun = session.stage === 'running' && session.ran !== null

	return (
		<div className="grid">
			<WorkflowStages
				copy={copy}
				format={format}
				human={choices.human}
				locked={running}
				onToggle={toggleHumanStage}
				run={
					showsRun
						? {
								playedCount: runPlayedCount,
								script,
							}
						: null
				}
			/>

			<CardSection className="mt-4 grid gap-4 pt-3 sm:grid-cols-2">
				<WorkflowChoice
					label={workflow.model}
					locked={running}
					name="model"
					note={format.message('workflow.modelPrice', {
						date: format.priceMonth,
						input: format.money(model.inputCents),
						output: format.money(model.outputCents),
					})}
					onChange={chooseModel}
					options={workflowModels.map((option) => ({
						label: option.name,
						value: option.id,
					}))}
					value={choices.model}
				/>
				<WorkflowChoice
					label={workflow.effort}
					locked={running}
					name="effort"
					onChange={chooseEffort}
					options={workflowEfforts.map((effort) => ({
						label: copy.analysis.efforts[effort],
						value: effort,
					}))}
					value={choices.effort}
				/>
			</CardSection>

			<CardSection className="mt-4 pt-3">
				<WorkflowSpend
					copy={copy}
					effort={choices.effort}
					figures={workflowFigures(choices)}
					format={format}
				/>
			</CardSection>

			<div className="mt-4 flex flex-wrap items-center justify-end gap-x-4 gap-y-2">
				{!canPlay && !running && (
					<span className="text-[11px] text-muted-foreground">
						{workflow.playUnavailable}
					</span>
				)}
				<button
					aria-disabled={!canPlay || undefined}
					className={playButton}
					onClick={handlePlay}
					type="button"
				>
					{playLabel(workflow, session, running)}
					{session.stage === 'launching' && (
						<span
							aria-hidden="true"
							className="inline-block size-3.5 rounded-full border-2 border-current border-r-transparent motion-safe:animate-spin"
						/>
					)}
				</button>
			</div>

			{showsRun && session.ran && (
				<CardSection className="mt-4 pt-3">
					<WorkflowResult
						copy={copy}
						figures={workflowFigures(session.ran)}
						format={format}
						playedCount={runPlayedCount}
						script={script}
					/>
				</CardSection>
			)}
		</div>
	)
}

export { type RunOutcome, runOutcome } from './run-outcome.compute'
export { workflowStatus } from './workflow-status.compute'
export { workflowSummary } from './workflow-summary.compute'
