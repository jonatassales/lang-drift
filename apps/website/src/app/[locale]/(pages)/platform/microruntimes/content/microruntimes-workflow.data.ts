import type { PlanEffort } from './microruntimes-content.types'
import type {
	HumanStageId,
	RunLine,
	Tokens,
	WorkflowChoices,
	WorkflowStageId,
} from './microruntimes-workflow.types'

// The hero's Workflow card and the run it plays: the Claude models the team
// can choose, their Anthropic list prices, the tokens a run of this ticket is
// expected to use at each effort, the budget the design doc sets, and the
// delivery stages. Money is integer cents; only the formatter divides. The
// run is simulated: its usage record, branch, hash and pull request are
// made up, like the artifact file names.

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

// runtime.md's own effort tokens; the API's levels above high are not
// offered.
export const workflowEfforts = [
	'low',
	'medium',
	'high',
] as const satisfies readonly PlanEffort[]

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

export const workflowStageIds = [
	'planning',
	'spec',
	'tickets',
	'implement',
	'dev-ticket',
	'review',
	'merge',
] as const

// The stages this workflow always gives a person: planning ends only when
// the owner confirms it, so its checkbox stays checked.
export const requiredHumanStages: readonly HumanStageId[] = ['planning']

// The skills each stage runs, as the repository names them.
export const stageSkills: Record<WorkflowStageId, string | null> = {
	planning: 'grilling',
	spec: 'to-spec',
	tickets: 'to-tickets',
	implement: 'implement · tdd',
	'dev-ticket': null,
	review: 'code-review',
	merge: null,
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

export const workflowDelivery = {
	branch: 'feat/0031-book-a-demo-runtime',
	hash: '3f9c2e1',
	pullRequest: 19,
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
	'devTicketWaiting',
	'devTicketClosed',
	'codeReview',
	'runtimeReady',
	'completed',
	'pushed',
	'pullRequest',
	'checks',
	'checksPassed',
	'safeSeatMet',
	'runtimeStopped',
	'readyForReview',
	'readyToMerge',
	'nextReview',
	'nextMerge',
] as const

// The lines a stage's own work prints once it is allowed to run.
export const stageWork: Record<
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
	],
}

export const workStageIds = [
	'planning',
	'spec',
	'tickets',
	'implement',
] as const
