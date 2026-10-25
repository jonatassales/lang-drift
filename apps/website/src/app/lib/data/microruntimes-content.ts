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

// The hero terminal's illustrative session with the connected agent: each
// line's kind, in the order it prints. The text lives in the catalog.
export const terminalSession = [
	{
		id: 'doctor',
		type: 'command',
	},
	{
		id: 'tools',
		type: 'success',
	},
	{
		id: 'run',
		type: 'command',
		delay: 600,
	},
	{
		id: 'sync',
		type: 'comment',
	},
	{
		id: 'runtime',
		type: 'output',
	},
	{
		id: 'configure',
		type: 'success',
	},
	{
		id: 'validate',
		type: 'success',
	},
	{
		id: 'prepare',
		type: 'success',
		delay: 700,
	},
	{
		id: 'execute',
		type: 'success',
		delay: 700,
	},
	{
		id: 'inspect',
		type: 'success',
	},
	{
		id: 'cleanup',
		type: 'success',
		delay: 700,
	},
	{
		id: 'evidence',
		type: 'output',
	},
] as const satisfies readonly {
	delay?: number
	id: string
	type: TerminalLineType
}[]

export type TerminalLineId = (typeof terminalSession)[number]['id']

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
		titleStart: string
		titleHighlight: string
		titleEnd: string
		lead: string
		primaryCta: string
		secondaryCta: string
		summary: string
	}
	lifecycle: {
		terminal: {
			title: string
			label: string
			lines: Record<TerminalLineId, string>
		}
		stepperLabel: string
		stages: Record<
			LifecycleStageId,
			{ label: string; does: string; leaves: string }
		>
		cpuValue: string
		memoryValue: string
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
