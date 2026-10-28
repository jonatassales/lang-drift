import { useAtomValue, useSetAtom } from 'jotai'

import {
	budgetCents,
	type RunLine,
	runHasPlayed,
	type WorkflowChoices,
	workflowDelivery,
	workflowEfforts,
	workflowFigures,
	workflowModel,
	workflowModels,
} from '@data/microruntimes-workflow'
import {
	canPlayWorkflowAtom,
	chooseEffortAtom,
	chooseModelAtom,
	type MicroruntimeSession,
	microruntimeSessionAtom,
	playWorkflowAtom,
	runInProgressAtom,
	runPlayedCountAtom,
	runScriptAtom,
	toggleHumanStageAtom,
	workflowChoicesAtom,
} from '@state'

import { CardSection } from './card-section'
import type { TicketCardCopy, TicketCardFormat } from './ticket-card-format'
import { WorkflowChoice } from './workflow-choice'
import { WorkflowResult } from './workflow-result'
import { WorkflowSpend } from './workflow-spend'
import { WorkflowStages } from './workflow-stages'

interface WorkflowCardProps {
	copy: TicketCardCopy
	format: TicketCardFormat
}

const playButton =
	'inline-flex min-h-[38px] cursor-pointer items-center gap-[7px] rounded-full bg-primary px-[17px] text-[13px] font-[650] text-primary-foreground motion-safe:transition-[transform,filter] motion-safe:duration-200 hover:saturate-[1.18] motion-safe:hover:-translate-y-px aria-disabled:cursor-not-allowed aria-disabled:opacity-60 aria-disabled:hover:translate-y-0'

// The label the Play button carries in each moment of the session.
function playLabel(
	workflow: TicketCardCopy['workflow'],
	session: MicroruntimeSession,
	running: boolean,
) {
	if (session.stage === 'launching') return workflow.starting
	if (running) return workflow.playing
	if (session.stage === 'running') return workflow.playAgain

	return workflow.play
}

// The session's last card: the run the plan leads to, summarized as money.
// The visitor picks the human stages, a Claude model and an effort, sees the
// expected spend against the budget the team set, and plays the workflow.
// The card then follows the run the terminal plays, from the gate to the
// pull request.
export function WorkflowCard(props: WorkflowCardProps) {
	const { copy, format } = props

	const { workflow } = copy
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
					onClick={() => play()}
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

interface RunOutcome {
	choices: WorkflowChoices
	playedCount: number
	script: readonly RunLine[]
}

// The run's end, when it has one: paused at the gate, or ready for the
// owner's review or to merge.
function runEnd(outcome: RunOutcome) {
	const { playedCount, script } = outcome
	const played = (message: Parameters<typeof runHasPlayed>[2]) =>
		runHasPlayed(script, playedCount, message)

	if (played('gatePaused')) return 'paused'
	if (played('readyForReview')) return 'readyForReview'
	if (played('readyToMerge')) return 'readyToMerge'

	return null
}

// The workflow's line in its collapsed header: the chosen model, effort and
// estimate against the budget, then how the last run ended.
export function workflowSummary(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	choices: WorkflowChoices,
	outcome: RunOutcome | null,
) {
	const end = outcome ? runEnd(outcome) : null

	if (end === 'paused') return copy.workflow.summaryPaused
	if (outcome && end) {
		return format.message('workflow.summaryDone', {
			number: workflowDelivery.pullRequest,
			observed: format.money(workflowFigures(outcome.choices).observed),
			state: copy.workflow[end],
		})
	}

	return format.message('workflow.summaryChoice', {
		budget: format.money(budgetCents),
		effort: copy.analysis.efforts[choices.effort],
		estimate: format.money(workflowFigures(choices).estimate),
		model: workflowModel(choices.model).shortName,
	})
}

// The one sentence the hero announces when a run ends, or nothing while it
// plays.
export function workflowStatus(
	copy: TicketCardCopy,
	format: TicketCardFormat,
	outcome: RunOutcome | null,
) {
	if (!outcome) return ''

	const end = runEnd(outcome)
	const figures = workflowFigures(outcome.choices)

	if (end === 'paused') {
		return format.message('workflow.statusPaused', {
			budget: format.money(budgetCents),
			reserved: format.money(figures.reserved),
		})
	}
	if (!end) return ''

	return format.message('workflow.statusDone', {
		number: workflowDelivery.pullRequest,
		observed: format.money(figures.observed),
		reserved: format.money(figures.reserved),
		state: copy.workflow[end],
	})
}

export function runOutcome(
	session: MicroruntimeSession,
	script: readonly RunLine[],
): RunOutcome | null {
	if (session.stage !== 'running' || !session.ran) return null

	return {
		choices: session.ran,
		playedCount: session.playedCount,
		script,
	}
}
