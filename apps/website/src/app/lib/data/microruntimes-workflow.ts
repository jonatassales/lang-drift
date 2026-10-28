import type { TerminalLineType } from '@langdrift/web/vendors/smoothui/mac-terminal'

import type { PlanEffort } from './microruntimes-content'

// The hero's Workflow card and the run it plays: the Claude models the team
// can choose, their Anthropic list prices, the tokens a run of this ticket is
// expected to use at each effort, the budget the design doc sets, and the
// delivery stages. Money is integer cents; only the formatter divides. The
// run is simulated: its usage record, branch, hash and pull request are
// illustrative data, like the artifact file names.

// The date the list prices below were read, for the card's price source.
export const priceDate = Date.UTC(2026, 9, 6)

export const workflowModels = [
	{
		id: 'claude-opus-5-5',
		name: 'Claude Opus 5.5',
		shortName: 'Opus 5.5',
		inputCents: 400,
		outputCents: 2000,
	},
	{
		id: 'claude-sonnet-5-5',
		name: 'Claude Sonnet 5.5',
		shortName: 'Sonnet 5.5',
		inputCents: 200,
		outputCents: 1000,
	},
	{
		id: 'claude-haiku-5-5',
		name: 'Claude Haiku 5.5',
		shortName: 'Haiku 5.5',
		inputCents: 10,
		outputCents: 50,
	},
] as const

export type WorkflowModel = (typeof workflowModels)[number]
export type WorkflowModelId = WorkflowModel['id']

// runtime.md's own effort tokens; the API's levels above high are not
// offered.
export const workflowEfforts = [
	'low',
	'medium',
	'high',
] as const satisfies readonly PlanEffort[]

export function workflowModel(id: WorkflowModelId): WorkflowModel {
	const model = workflowModels.find((candidate) => candidate.id === id)
	if (!model) throw new Error(`Unknown workflow model: ${id}`)

	return model
}

// Input and output tokens, in thousands.
export interface Tokens {
	inputThousands: number
	outputThousands: number
}

// What a run of this ticket is expected to use, from the spec's five
// components; no episode measured tokens.
export const expectedTokens: Record<PlanEffort, Tokens> = {
	low: {
		inputThousands: 2000,
		outputThousands: 120,
	},
	medium: {
		inputThousands: 3500,
		outputThousands: 250,
	},
	high: {
		inputThousands: 6000,
		outputThousands: 500,
	},
}

// The simulated usage record a run leaves, a little above the expectation.
export const observedTokens: Record<PlanEffort, Tokens> = {
	low: {
		inputThousands: 2100,
		outputThousands: 125,
	},
	medium: {
		inputThousands: 3700,
		outputThousands: 240,
	},
	high: {
		inputThousands: 6300,
		outputThousands: 520,
	},
}

// The budget the design doc sets for a run's model calls, and the margin
// the gate reserves above the estimate.
export const budgetCents = 2500
export const reserveMarginPercent = 20

// What the tokens cost at the model's uncached list prices.
export function spendCents(modelId: WorkflowModelId, tokens: Tokens) {
	const model = workflowModel(modelId)
	const tenthsOfCents =
		tokens.inputThousands * model.inputCents +
		tokens.outputThousands * model.outputCents

	return Math.round(tenthsOfCents / 1000)
}

export function reserveCents(estimateCents: number) {
	return Math.ceil((estimateCents * (100 + reserveMarginPercent)) / 100)
}

export const workflowStageIds = [
	'planning',
	'spec',
	'tickets',
	'implement',
	'review',
	'merge',
] as const

export type WorkflowStageId = (typeof workflowStageIds)[number]

// The stages a person can be asked to confirm; merge is always the owner's.
export type HumanStageId = Exclude<WorkflowStageId, 'merge'>

// The stages this workflow always gives a person: planning ends only when
// the owner confirms it, so its checkbox stays checked.
const requiredHumanStages: readonly HumanStageId[] = ['planning']

export function isRequiredHumanStage(stage: HumanStageId) {
	return requiredHumanStages.includes(stage)
}

// The skills each stage runs, as the repository names them.
export const stageSkills: Record<WorkflowStageId, string | null> = {
	planning: 'grilling',
	spec: 'to-spec',
	tickets: 'to-tickets',
	implement: 'implement · tdd · code-review',
	review: null,
	merge: null,
}

export interface WorkflowChoices {
	effort: PlanEffort
	human: Record<HumanStageId, boolean>
	model: WorkflowModelId
}

export const initialWorkflowChoices: WorkflowChoices = {
	effort: 'medium',
	human: {
		planning: true,
		spec: false,
		tickets: false,
		implement: false,
		review: true,
	},
	model: 'claude-sonnet-5-5',
}

// The run's money in cents: the estimate, the reservation the gate would
// hold, what the usage record observed and what the reservation releases.
export function workflowFigures(choices: WorkflowChoices) {
	const estimate = spendCents(choices.model, expectedTokens[choices.effort])
	const reserved = reserveCents(estimate)
	const observed = spendCents(choices.model, observedTokens[choices.effort])

	return {
		estimate,
		reserved,
		observed,
		released: reserved - observed,
		overBudget: reserved > budgetCents,
	}
}

export type WorkflowFigures = ReturnType<typeof workflowFigures>

export const workflowDelivery = {
	branch: 'feat/0031-early-access-runtime',
	hash: '3f9c2e1',
	pullRequest: 19,
}

// US dollars with cents, for figures from a fraction of a cent upwards.
export function workflowDollars(locale: string) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 2,
		minimumFractionDigits: 2,
		style: 'currency',
	})
}

// Thousands of tokens as a compact count: 3.5M, 250K.
export function tokenCount(locale: string, thousands: number) {
	return new Intl.NumberFormat(locale, {
		maximumFractionDigits: 1,
		notation: 'compact',
	}).format(thousands * 1000)
}

// The run's lines in the catalog, by message. Waiting and confirmed lines
// repeat per stage, so a line carries its own unique id beside the message.
export const runMessageIds = [
	'command',
	'scope',
	'intent',
	'estimate',
	'budget',
	'gateOpened',
	'gatePaused',
	'paused',
	'waiting',
	'confirmed',
	'grilling',
	'toSpec',
	'toTickets',
	'hook',
	'implement',
	'codeReview',
	'runtimeReady',
	'completed',
	'pushed',
	'pullRequest',
	'checks',
	'checksPassed',
	'runtimeStopped',
	'readyForReview',
	'readyToMerge',
	'nextReview',
	'nextMerge',
] as const

export type RunMessageId = (typeof runMessageIds)[number]

export interface RunLine {
	delay: number
	id: string
	message: RunMessageId
	stage?: HumanStageId
	type: TerminalLineType
	values?: Record<string, string>
}

// The lines a stage's own work prints once it is allowed to run.
const stageWork: Record<
	Exclude<HumanStageId, 'review'>,
	readonly Omit<RunLine, 'id' | 'stage'>[]
> = {
	planning: [
		{
			delay: 800,
			message: 'grilling',
			type: 'output',
		},
	],
	spec: [
		{
			delay: 700,
			message: 'toSpec',
			type: 'output',
		},
	],
	tickets: [
		{
			delay: 700,
			message: 'toTickets',
			type: 'output',
		},
	],
	implement: [
		{
			delay: 500,
			message: 'hook',
			type: 'output',
		},
		{
			delay: 1200,
			message: 'implement',
			type: 'output',
		},
		{
			delay: 800,
			message: 'codeReview',
			type: 'output',
		},
	],
}

const workStageIds = ['planning', 'spec', 'tickets', 'implement'] as const

// One stage: the wait for the owner when the stage asks for a person, then
// the stage's own work.
function stageLines(stage: (typeof workStageIds)[number], human: boolean) {
	const confirmation: RunLine[] = human
		? [
				{
					delay: 600,
					id: `waiting-${stage}`,
					message: 'waiting',
					stage,
					type: 'output',
					values: { stage },
				},
				{
					delay: 1600,
					id: `confirmed-${stage}`,
					message: 'confirmed',
					stage,
					type: 'success',
					values: { stage },
				},
			]
		: []
	const work = stageWork[stage].map((line) => ({
		...line,
		id: line.message,
		stage,
	}))

	return [...confirmation, ...work]
}

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

// True once the run has completed a line with this message.
export function runHasPlayed(
	script: readonly RunLine[],
	playedCount: number,
	message: RunMessageId,
) {
	const index = script.findIndex((line) => line.message === message)

	return index >= 0 && playedCount > index
}

export type StageProgress = 'pending' | 'active' | 'done'

// Where a stage stands in a playing run: pending until its first line
// prints, active while its lines print, done after its last.
export function stageProgress(
	script: readonly RunLine[],
	playedCount: number,
	stage: HumanStageId,
): StageProgress {
	const indexes = script
		.map((line, index) => (line.stage === stage ? index : -1))
		.filter((index) => index >= 0)
	const first = indexes.at(0)
	const last = indexes.at(-1)

	if (first === undefined || last === undefined) return 'pending'
	if (playedCount > last) return 'done'
	if (playedCount > first) return 'active'

	return 'pending'
}
