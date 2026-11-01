import {
	budgetCents,
	expectedTokens,
	type RunLine,
	type WorkflowChoices,
	workflowDelivery,
	workStageIds,
} from './microruntimes-workflow'
import { stageLines } from './stage-lines.compute'
import { tokenCount } from './token-count.fmt'
import { workflowDollars } from './workflow-dollars.fmt'
import { workflowFigures } from './workflow-figures.compute'

// The run the Workflow card plays, built from the choices it was launched
// with: the recorded intent, the estimate and the budget, the console's gate,
// then each stage, the microruntime the checks use, the push, the pull
// request and its checks. A reservation over the budget stops the run at
// the gate with nothing paid for. Dollar figures are formatted in English,
// as every console line is.
export function runSession(choices: WorkflowChoices): RunLine[] {
	const dollars = workflowDollars('en')
	const figures = workflowFigures(choices)
	const tokens = expectedTokens[choices.effort]
	const money = (cents: number) => dollars.format(cents / 100)
	const opening: RunLine[] = [
		{
			delay: 0,
			id: 'command',
			message: 'command',
			type: 'command',
		},
		{
			delay: 400,
			id: 'scope',
			message: 'scope',
			type: 'comment',
		},
		{
			delay: 600,
			id: 'intent',
			message: 'intent',
			type: 'output',
			values: {
				effort: choices.effort,
				model: choices.model,
			},
		},
		{
			delay: 600,
			id: 'estimate',
			message: 'estimate',
			type: 'output',
			values: {
				estimate: money(figures.estimate),
				input: tokenCount('en', tokens.inputThousands),
				output: tokenCount('en', tokens.outputThousands),
			},
		},
		{
			delay: 500,
			id: 'budget',
			message: 'budget',
			type: 'output',
			values: { budget: money(budgetCents) },
		},
	]

	if (figures.overBudget) {
		return [
			...opening,
			{
				delay: 900,
				id: 'gatePaused',
				message: 'gatePaused',
				type: 'warning',
				values: {
					budget: money(budgetCents),
					reserved: money(figures.reserved),
				},
			},
			{
				delay: 700,
				id: 'paused',
				message: 'paused',
				type: 'comment',
			},
		]
	}

	const reviewed = choices.human.review
	const delivery = {
		branch: workflowDelivery.branch,
		hash: workflowDelivery.hash,
		number: String(workflowDelivery.pullRequest),
	}

	return [
		...opening,
		{
			delay: 900,
			id: 'gateOpened',
			message: 'gateOpened',
			type: 'success',
			values: {
				budget: money(budgetCents),
				reserved: money(figures.reserved),
			},
		},
		...workStageIds.flatMap((stage) => stageLines(stage, choices.human[stage])),
		{
			delay: 900,
			id: 'runtimeReady',
			message: 'runtimeReady',
			type: 'output',
		},
		{
			delay: 800,
			id: 'completed',
			message: 'completed',
			type: 'success',
			values: {
				observed: money(figures.observed),
				reserved: money(figures.reserved),
			},
		},
		{
			delay: 600,
			id: 'pushed',
			message: 'pushed',
			type: 'success',
			values: delivery,
		},
		{
			delay: 600,
			id: 'pullRequest',
			message: 'pullRequest',
			type: 'success',
			values: delivery,
		},
		{
			delay: 700,
			id: 'checks',
			message: 'checks',
			type: 'output',
		},
		{
			delay: 1800,
			id: 'checksPassed',
			message: 'checksPassed',
			type: 'success',
		},
		{
			delay: 600,
			id: 'runtimeStopped',
			message: 'runtimeStopped',
			type: 'output',
		},
		{
			delay: 800,
			id: 'ready',
			message: reviewed ? 'readyForReview' : 'readyToMerge',
			stage: 'review',
			type: 'success',
			values: {
				...delivery,
				released: money(figures.released),
			},
		},
		{
			delay: 700,
			id: 'next',
			message: reviewed ? 'nextReview' : 'nextMerge',
			type: 'comment',
			values: delivery,
		},
	]
}
