import type { TerminalLineType } from '@langdrift/web/vendors/smoothui/mac-terminal'

// Illustrative Microruntimes content: the ids, limits, identities and times
// the Microruntimes page and its demos render. Sentences live in the
// microruntimes catalog. The runtime file mirrors the contract of the ex
// connected agent (one langdrift JSON fence in runtime.md); its values are
// simulated.

export const lifecycleStageIds = [
	'configure',
	'validate',
	'prepare',
	'execute',
	'inspect',
	'cleanup',
] as const

export type LifecycleStageId = (typeof lifecycleStageIds)[number]

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
	'availability',
] as const

export type LifecyclePhaseId = (typeof lifecyclePhaseIds)[number]
export type LifecycleScenarioId = (typeof lifecycleScenarioIds)[number]
export type DescriptorAnnotationId = (typeof descriptorAnnotationIds)[number]
export type OwnershipRowId = (typeof ownershipRowIds)[number]
export type StateCardId = (typeof stateCardIds)[number]
export type FlowTargetId = (typeof flowTargetIds)[number]
export type FlowGuaranteeId = (typeof flowGuaranteeIds)[number]
export type MicroruntimesFaqId = (typeof microruntimesFaqIds)[number]

// The hero terminal's simulated planning session: each line's kind and the
// wait before it prints, in order. The text lives in the catalog. The typed
// command is the intended entry point of the connected agent; the planning
// lines describe the console's proposed flow, not a shipped planner.
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
		id: 'memory',
		type: 'output',
		delay: 900,
	},
	{
		id: 'episodes',
		type: 'output',
		delay: 1400,
	},
	{
		id: 'suggest',
		type: 'success',
		delay: 900,
	},
	{
		id: 'reads',
		type: 'output',
		delay: 700,
	},
	{
		id: 'guardrail',
		type: 'output',
		delay: 600,
	},
	{
		id: 'bounds',
		type: 'output',
		delay: 900,
	},
	{
		id: 'resolve',
		type: 'output',
		delay: 800,
	},
	{
		id: 'go',
		type: 'success',
		delay: 650,
	},
	{
		id: 'node',
		type: 'success',
		delay: 650,
	},
	{
		id: 'nats',
		type: 'success',
		delay: 650,
	},
	{
		id: 'redis',
		type: 'success',
		delay: 650,
	},
	{
		id: 'next',
		type: 'success',
		delay: 650,
	},
	{
		id: 'size',
		type: 'success',
		delay: 1200,
	},
	{
		id: 'limits',
		type: 'warning',
		delay: 700,
	},
	{
		id: 'price',
		type: 'output',
		delay: 800,
	},
	{
		id: 'draft',
		type: 'success',
		delay: 1200,
	},
	{
		id: 'handoff',
		type: 'comment',
		delay: 800,
	},
] as const satisfies readonly {
	delay?: number
	id: string
	type: TerminalLineType
}[]

export type PlanLineId = (typeof planSession)[number]['id']

export function planLineIndex(id: PlanLineId) {
	return planSession.findIndex((line) => line.id === id)
}

// True once the terminal has completed the line, so the card can show what
// that line reveals.
export function hasPlayed(playedCount: number, id: PlanLineId) {
	return playedCount > planLineIndex(id)
}

export type PlanEffort = 'low' | 'medium' | 'high'
export type ComponentKind = 'api' | 'broker' | 'cache' | 'web'
export type Technology = 'Go' | 'Node.js' | 'NATS' | 'Redis' | 'Next.js'

// A CPU and memory pair, in the units runtime.md bounds: millicores and MiB.
export interface Resources {
	cpuMillis: number
	memoryMiB: number
}

// One component the plan resolves from the spec: where it lives, what it
// is, and the share of the envelope the plan proposes for it. The id is the
// session line that resolves it.
export interface PlanComponent extends Resources {
	id: PlanLineId
	kind: ComponentKind
	path: string
	technology: Technology
}

// The ticket the session plans, as the agent reads it from the project's
// artifacts: the numbers follow the scaffold ex init produces. The envelope
// is the team's, from the design doc; the components' shares are the plan's
// proposal and must stay inside it.
export const ticketPlan = {
	ticket: {
		number: '0031',
		slug: 'early-access-runtime',
	},
	artifacts: {
		prd: '0002-PRD-early-access.prd.md',
		designDoc: '0005-DD-early-access-runtime.dd.md',
		spec: '0011-SPEC-early-access-runtime.spec.md',
	},
	guardrail: {
		file: '0001-finops.gr.md',
		criteria: ['FIN-02', 'FIN-04'],
	},
	model: 'operator/allowed-model',
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
			memoryMiB: 256,
			path: 'services/signup-api',
			technology: 'Go',
		},
		{
			cpuMillis: 300,
			id: 'node',
			kind: 'api',
			memoryMiB: 512,
			path: 'services/notify-api',
			technology: 'Node.js',
		},
		{
			cpuMillis: 200,
			id: 'nats',
			kind: 'broker',
			memoryMiB: 256,
			path: 'platform/nats',
			technology: 'NATS',
		},
		{
			cpuMillis: 200,
			id: 'redis',
			kind: 'cache',
			memoryMiB: 256,
			path: 'platform/redis',
			technology: 'Redis',
		},
		{
			cpuMillis: 500,
			id: 'next',
			kind: 'web',
			memoryMiB: 1024,
			path: 'apps/signup',
			technology: 'Next.js',
		},
	] satisfies readonly PlanComponent[],
	draftRevision: 1,
} as const

// The totals the plan proposes: the sum of the components' shares.
export function suggestedTotals(
	components: readonly PlanComponent[],
): Resources {
	return components.reduce<Resources>(
		(totals, component) => ({
			cpuMillis: totals.cpuMillis + component.cpuMillis,
			memoryMiB: totals.memoryMiB + component.memoryMiB,
		}),
		{
			cpuMillis: 0,
			memoryMiB: 0,
		},
	)
}

// Which line of the session reveals which part of the card. Each component
// row is revealed by its own line.
export const planReveals = {
	card: 'plan',
	suggestionSkeleton: 'memory',
	suggestion: 'suggest',
	artifactsRead: 'reads',
	guardrail: 'guardrail',
	componentsSkeleton: 'resolve',
	bounded: 'bounds',
	suggested: 'size',
	limits: 'limits',
	priceBasis: 'price',
	report: 'draft',
	finished: 'handoff',
} as const satisfies Record<string, PlanLineId>

export type PlanRevealId = keyof typeof planReveals

// True once the terminal has completed the line that reveals this part of
// the card.
export function revealed(playedCount: number, id: PlanRevealId) {
	return hasPlayed(playedCount, planReveals[id])
}

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
	model: 'operator/allowed-model',
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

// Renders the descriptor as the indented JSON a reader would see in the
// file, with short arrays kept on one line.
export function formatDescriptor(descriptor: typeof runtimeDescriptor) {
	return JSON.stringify(descriptor, null, 2).replace(
		/\[\n\s+([^\]{}]*?)\n\s+\]/g,
		(_match, items: string) => `[${items.replace(/,\n\s+/g, ', ')}]`,
	)
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

export interface TitledCopy {
	title: string
	description: string
}

export interface MicroruntimesMessages {
	skip: string
	hero: {
		title: string
		summary: string
	}
	session: {
		terminal: {
			title: string
			label: string
			lines: Record<PlanLineId, string>
		}
		card: {
			ticket: string
			status: string
			title: string
			readMark: string
			notSet: string
			resolving: string
			artifacts: Record<'prd' | 'designDoc' | 'spec' | 'guardrail', string>
			guardrailNote: string
			suggestedModel: string
			suggestedEffort: string
			efforts: Record<PlanEffort, string>
			kinds: Record<ComponentKind, string>
			recommendationNote: string
			components: string
			columns: { cpu: string; memory: string }
			resolvedProgress: string
			componentsTotal: string
			figures: string
			hardware: Record<
				| 'resource'
				| 'bounded'
				| 'suggested'
				| 'cpu'
				| 'memory'
				| 'gpu'
				| 'disk',
				string
			>
			notRuntimeField: string
			none: string
			priceBasis: string
			report: Record<
				'bounded' | 'suggested' | 'suggestedResource' | 'final',
				string
			>
			deployments: string
			pendingAuthorization: string
			finalDetail: string
			planPath: string
			control: { pause: string; resume: string; replay: string }
			cpuValue: string
			memoryValue: string
		}
	}
	lifecycle: {
		stepperLabel: string
		stages: Record<
			LifecycleStageId,
			{ label: string; does: string; leaves: string }
		>
		whatItDoes: string
		evidenceItLeaves: string
	}
	scenarios: {
		title: string
		lead: string
		selectorLabel: string
		phases: Record<LifecyclePhaseId, string>
		items: Record<
			LifecycleScenarioId,
			{ tabLabel: string; context: string } & Record<LifecyclePhaseId, string>
		>
		outro: string
		outroLink: string
	}
	descriptor: {
		kicker: string
		title: string
		lead: string
		annotations: Record<DescriptorAnnotationId, TitledCopy>
	}
	ownership: {
		kicker: string
		title: string
		lead: string
		concern: string
		consumer: string
		platform: string
		rows: Record<
			OwnershipRowId,
			{ concern: string; consumer: string; platform: string }
		>
	}
	state: {
		kicker: string
		title: string
		lead: string
		cards: Record<StateCardId, TitledCopy & { eyebrow: string }>
	}
	evidence: {
		kicker: string
		title: string
		lead: string
		ledgerTitle: string
		rows: {
			execution: string
			started: string
			stopped: string
			owned: string
			removed: string
			leftovers: string
			evidence: string
		}
		executionValue: string
		resources: string
		noLeftovers: string
		selectedStage: string
		stages: Record<LifecycleStageId, { kind: string; title: string }>
	}
	flow: {
		kicker: string
		title: string
		lead: string
		console: TitledCopy
		agent: TitledCopy
		targetsTitle: string
		targets: Record<FlowTargetId, TitledCopy>
		guarantees: Record<FlowGuaranteeId, TitledCopy>
	}
	faq: {
		title: string
		description: string
		items: Record<MicroruntimesFaqId, { question: string; answer: string }>
	}
}
