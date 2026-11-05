import type {
	DocumentId,
	Episode,
	LifecycleStageId,
	PlanComponent,
	PlanEffort,
	PlanLineId,
	PlanPhaseId,
	PlanSessionLine,
	Resources,
} from './microruntimes-content.types'
import type { WorkflowModelId } from './microruntimes-workflow.types'

// Microruntimes content: the ids, limits, identities and times the
// Microruntimes page and its demos render. Sentences live in the
// microruntimes catalog. The runtime file mirrors the contract of the
// connected agent of the LangDrift CLI (one langdrift JSON fence in
// runtime.md); its values are simulated.

export const lifecycleStageIds = [
	'configure',
	'validate',
	'prepare',
	'execute',
	'inspect',
	'cleanup',
] as const

export const initialStageId: LifecycleStageId = 'execute'

export const lifecyclePhaseIds = [
	'compose',
	'run',
	'inspect',
	'cleanup',
] as const
export const lifecycleScenarioIds = [
	'add-component',
	'preparation-fails',
	'lease-lost',
] as const
export const descriptorAnnotationIds = [
	'identity',
	'components',
	'build',
	'resources',
	'reasoning',
] as const
export const ownershipRowIds = [
	'source',
	'composition',
	'build',
	'lifecycle',
	'cleanup',
	'evidence',
] as const
export const stateCardIds = ['requested', 'observed', 'unknown'] as const
export const flowTargetIds = ['kind', 'terraform', 'kubernetes'] as const
export const flowGuaranteeIds = ['claims', 'events', 'ownership'] as const
export const microruntimesFaqIds = [
	'what',
	'source',
	'stacks',
	'sandbox',
	'license',
	'availability',
] as const

// The hero terminal's simulated planning session: each line's kind and the
// wait before it prints, in order. The text lives in the catalog. The typed
// command is the intended entry point of the connected agent; the planning
// lines describe the console's proposed flow, not a shipped planner. The
// session plays in three acts, each opened by a comment: the ticket, the
// analysis of its documents and the project's episodes, and the plan.
export const planSession = [
	{
		id: 'plan',
		type: 'command',
		delay: 0,
	},
	{
		id: 'scope',
		type: 'comment',
		delay: 400,
	},
	{
		id: 'readTicket',
		type: 'output',
		delay: 500,
	},
	{
		id: 'links',
		type: 'output',
		delay: 600,
	},
	{
		id: 'analysis',
		type: 'comment',
		delay: 900,
	},
	{
		id: 'readPrd',
		type: 'output',
		delay: 500,
	},
	{
		id: 'readDesignDoc',
		type: 'output',
		delay: 450,
	},
	{
		id: 'readSpec',
		type: 'output',
		delay: 450,
	},
	{
		id: 'reads',
		type: 'success',
		delay: 600,
	},
	{
		id: 'guardrail',
		type: 'output',
		delay: 500,
	},
	{
		id: 'judgeGuardrail',
		type: 'output',
		delay: 900,
	},
	{
		id: 'bounds',
		type: 'output',
		delay: 700,
	},
	{
		id: 'budget',
		type: 'output',
		delay: 700,
	},
	{
		id: 'memory',
		type: 'output',
		delay: 800,
	},
	{
		id: 'embed',
		type: 'output',
		delay: 450,
	},
	{
		id: 'vectorSearch',
		type: 'output',
		delay: 500,
	},
	{
		id: 'candidates',
		type: 'output',
		delay: 600,
	},
	{
		id: 'judgeRelevance',
		type: 'output',
		delay: 1000,
	},
	{
		id: 'drop',
		type: 'output',
		delay: 700,
	},
	{
		id: 'episodes',
		type: 'output',
		delay: 600,
	},
	{
		id: 'reasonResources',
		type: 'output',
		delay: 1000,
	},
	{
		id: 'weighEffort',
		type: 'output',
		delay: 900,
	},
	{
		id: 'weighModel',
		type: 'output',
		delay: 900,
	},
	{
		id: 'suggest',
		type: 'success',
		delay: 800,
	},
	{
		id: 'planning',
		type: 'comment',
		delay: 900,
	},
	{
		id: 'resolve',
		type: 'output',
		delay: 600,
	},
	{
		id: 'inspectGo',
		type: 'output',
		delay: 400,
	},
	{
		id: 'go',
		type: 'success',
		delay: 450,
	},
	{
		id: 'inspectNode',
		type: 'output',
		delay: 400,
	},
	{
		id: 'node',
		type: 'success',
		delay: 450,
	},
	{
		id: 'inspectNats',
		type: 'output',
		delay: 400,
	},
	{
		id: 'nats',
		type: 'success',
		delay: 450,
	},
	{
		id: 'inspectRedis',
		type: 'output',
		delay: 400,
	},
	{
		id: 'redis',
		type: 'success',
		delay: 450,
	},
	{
		id: 'inspectNext',
		type: 'output',
		delay: 400,
	},
	{
		id: 'next',
		type: 'success',
		delay: 450,
	},
	{
		id: 'crossCheck',
		type: 'output',
		delay: 800,
	},
	{
		id: 'size',
		type: 'success',
		delay: 900,
	},
	{
		id: 'limits',
		type: 'warning',
		delay: 700,
	},
	{
		id: 'price',
		type: 'output',
		delay: 700,
	},
	{
		id: 'validate',
		type: 'output',
		delay: 700,
	},
	{
		id: 'draft',
		type: 'success',
		delay: 1000,
	},
	{
		id: 'handoff',
		type: 'comment',
		delay: 2500,
	},
] as const satisfies readonly PlanSessionLine[]

// The acts of the session, each the concern of one card beside the
// terminal, and the line that opens each. The workflow opens once the plan
// is handed to the team, and its card plays the run.
export const planPhaseIds = ['ticket', 'analysis', 'plan', 'workflow'] as const

export const planPhases = {
	ticket: 'plan',
	analysis: 'analysis',
	plan: 'planning',
	workflow: 'handoff',
} as const satisfies Record<PlanPhaseId, PlanLineId>

// The documents the ticket links, in the order the agent reads them, and
// the session line that reads each.
export const documentIds = ['prd', 'designDoc', 'spec', 'guardrail'] as const

export const documentReads = {
	prd: 'readPrd',
	designDoc: 'readDesignDoc',
	spec: 'readSpec',
	guardrail: 'guardrail',
} as const satisfies Record<DocumentId, PlanLineId>

// The ticket the session plans, as the agent reads it from the project's
// artifacts: the numbers follow the scaffold dft init produces. The envelope
// is the team's, from the design doc; the components' shares are the plan's
// proposal and must stay inside it.
export const ticketPlan = {
	ticket: {
		number: '0031',
		slug: 'book-a-demo-runtime',
	},
	artifacts: {
		prd: '0002-PRD-book-a-demo.prd.md',
		designDoc: '0005-DD-book-a-demo-runtime.dd.md',
		spec: '0011-SPEC-book-a-demo-runtime.spec.md',
		guardrail: '0001-finops.gr.md',
	} satisfies Record<DocumentId, string>,
	guardrailCriteria: ['FIN-02', 'FIN-03', 'FIN-04'],
	spec: {
		components: 5,
		acceptanceCriteria: 3,
	},
	model: 'claude-sonnet-5-5' satisfies WorkflowModelId,
	effort: 'medium' satisfies PlanEffort,
	bounded: {
		cpuMillis: 2000,
		memoryMiB: 3072,
	} satisfies Resources,
	components: [
		{
			cpuMillis: 250,
			id: 'go',
			kind: 'api',
			manifest: 'go.mod',
			memoryMiB: 256,
			path: 'services/signup-api',
			port: 8080,
			technology: 'Go',
		},
		{
			cpuMillis: 300,
			id: 'node',
			kind: 'api',
			manifest: 'package.json',
			memoryMiB: 512,
			path: 'services/notify-api',
			port: 8081,
			technology: 'Node.js',
		},
		{
			cpuMillis: 200,
			id: 'nats',
			kind: 'broker',
			manifest: 'nats.conf',
			memoryMiB: 256,
			path: 'platform/nats',
			port: 4222,
			technology: 'NATS',
		},
		{
			cpuMillis: 200,
			id: 'redis',
			kind: 'cache',
			manifest: 'redis.conf',
			memoryMiB: 256,
			path: 'platform/redis',
			port: 6379,
			technology: 'Redis',
		},
		{
			cpuMillis: 500,
			id: 'next',
			kind: 'web',
			manifest: 'next.config.ts',
			memoryMiB: 1024,
			path: 'apps/signup',
			port: 3000,
			technology: 'Next.js',
		},
	] satisfies readonly PlanComponent[],
	draftRevision: 1,
} as const

// Below this similarity a candidate is dropped before it is cited.
export const similarityThreshold = 0.55

export const episodes: readonly Episode[] = [
	{
		id: 'ep-0412',
		similarity: 0.81,
		slug: 'checkout-preview-sizing',
		verdict: 'kept',
	},
	{
		id: 'ep-0388',
		similarity: 0.74,
		slug: 'notify-api-rollout',
		verdict: 'kept',
	},
	{
		id: 'ep-0351',
		similarity: 0.66,
		slug: 'signup-redesign-preview',
		verdict: 'kept',
	},
	{
		id: 'ep-0297',
		reason: 'docsOnly',
		similarity: 0.52,
		slug: 'docs-search-index',
		verdict: 'dropped',
	},
	{
		id: 'ep-0240',
		reason: 'otherRuntime',
		similarity: 0.47,
		slug: 'batch-export-worker',
		verdict: 'dropped',
	},
	{
		id: 'ep-0188',
		reason: 'noBroker',
		similarity: 0.42,
		slug: 'storefront-cache-warmup',
		verdict: 'dropped',
	},
]

// Which line of the session reveals which part of the cards. Each component
// row is revealed by its own line, each document by its read.
export const planReveals = {
	cards: 'plan',
	links: 'links',
	documentsRead: 'reads',
	guardrailJudged: 'judgeGuardrail',
	bounded: 'bounds',
	suggestionSkeleton: 'memory',
	episodesSkeleton: 'vectorSearch',
	candidates: 'candidates',
	verdicts: 'drop',
	cited: 'episodes',
	measured: 'reasonResources',
	suggestion: 'suggest',
	componentsSkeleton: 'resolve',
	crossChecked: 'crossCheck',
	suggested: 'size',
	limits: 'limits',
	priceBasis: 'price',
	validated: 'validate',
	drafted: 'draft',
	finished: 'handoff',
} as const satisfies Record<string, PlanLineId>

export const runtimeFilePath = '.langdrift/runtimes/checkout-preview/runtime.md'

// The JSON inside the runtime file's one langdrift fence, as the console
// would author it.
export const runtimeDescriptor = {
	version: 1,
	runtimeUid: 'rt-7f3c2a',
	runtimeId: 'checkout-preview',
	projectId: 'shop',
	revision: 3,
	name: 'Checkout preview',
	cpuMillis: 1500,
	memoryMiB: 2048,
	model: 'claude-sonnet-5-5',
	effort: 'medium',
	components: [
		{
			path: 'apps/storefront',
			image: 'storefront:local',
			port: 3000,
			readinessPath: '/healthz',
			build: {
				executable: 'docker',
				args: ['build', '-t', 'storefront:local', '.'],
			},
		},
		{
			path: 'services/checkout',
			image: 'checkout:local',
			port: 8080,
			readinessPath: '/ready',
			build: {
				executable: 'pnpm',
				args: ['run', 'build:image'],
			},
		},
		{
			path: 'services/payments-stub',
			image: 'payments-stub:local',
			port: 8081,
			readinessPath: '/ready',
			build: {
				executable: './scripts/build-image.sh',
				args: [],
			},
		},
	],
}

export const executionLedger = {
	executionId: 'exec-41b2',
	revision: 3,
	startedAt: Date.UTC(2026, 9, 24, 14, 2),
	stoppedAt: Date.UTC(2026, 9, 24, 15, 37),
	ownedResources: 7,
	removedResources: 7,
	evidencePath: '.audits/evidence/runs/exec-41b2/',
}

export const stageEvidencePaths: Record<LifecycleStageId, string> = {
	configure: '.langdrift/runtimes/checkout-preview/runtime.md#revision-3',
	validate: '.audits/evidence/runs/exec-41b2/validate.json',
	prepare: '.audits/evidence/runs/exec-41b2/prepare.json',
	execute: '.audits/evidence/runs/exec-41b2/claim.json',
	inspect: '.audits/evidence/runs/exec-41b2/inspect.json',
	cleanup: '.audits/evidence/runs/exec-41b2/cleanup.json',
}
