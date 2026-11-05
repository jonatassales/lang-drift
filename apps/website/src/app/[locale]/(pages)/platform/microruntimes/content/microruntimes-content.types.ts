import type { FaqItemCopy } from '@routes/platform/product-page/faq-item-copy.type'
import type { TerminalLineType } from 'web/vendors/smoothui/mac-terminal'

import type {
	descriptorAnnotationIds,
	documentIds,
	flowGuaranteeIds,
	flowTargetIds,
	lifecyclePhaseIds,
	lifecycleScenarioIds,
	lifecycleStageIds,
	microruntimesFaqIds,
	ownershipRowIds,
	planPhaseIds,
	planReveals,
	planSession,
	stateCardIds,
} from './microruntimes-content.data'
import type {
	RunMessageId,
	WorkflowStageId,
} from './microruntimes-workflow.types'

export type LifecycleStageId = (typeof lifecycleStageIds)[number]

export type LifecyclePhaseId = (typeof lifecyclePhaseIds)[number]

export type LifecycleScenarioId = (typeof lifecycleScenarioIds)[number]

export type DescriptorAnnotationId = (typeof descriptorAnnotationIds)[number]

export type OwnershipRowId = (typeof ownershipRowIds)[number]

export type StateCardId = (typeof stateCardIds)[number]

export type FlowTargetId = (typeof flowTargetIds)[number]

export type FlowGuaranteeId = (typeof flowGuaranteeIds)[number]

export type MicroruntimesFaqId = (typeof microruntimesFaqIds)[number]

export interface PlanSessionLine {
	delay?: number
	id: string
	type: TerminalLineType
}

export type PlanLineId = (typeof planSession)[number]['id']

export type PlanPhaseId = (typeof planPhaseIds)[number]

export type PlanEffort = 'low' | 'medium' | 'high'

export type ComponentKind = 'api' | 'broker' | 'cache' | 'web'

export type Technology = 'Go' | 'Node.js' | 'NATS' | 'Redis' | 'Next.js'

// A CPU and memory pair, in the units runtime.md bounds: millicores and MiB.
export interface Resources {
	cpuMillis: number
	memoryMiB: number
}

// One component the plan resolves from the spec: where it lives, what it
// is, the manifest and port the agent inspects, and the share of the
// envelope the plan proposes for it. The id is the session line that
// resolves it.
export interface PlanComponent extends Resources {
	id: PlanLineId
	kind: ComponentKind
	manifest: string
	path: string
	port: number
	technology: Technology
}

export type DocumentId = (typeof documentIds)[number]

export type EpisodeVerdict = 'kept' | 'dropped'

export type DropReason = 'otherRuntime' | 'noBroker' | 'docsOnly'

// One of the project's episodes the retrieval ranks for the ticket: its
// similarity to the ticket's intent, and whether the judgment kept it as a
// precedent or dropped it, with the reason. The slugs are made up.
export interface Episode {
	id: string
	reason?: DropReason
	similarity: number
	slug: string
	verdict: EpisodeVerdict
}

export type PlanRevealId = keyof typeof planReveals

export interface TitledCopy {
	title: string
	description: string
}

export interface LifecycleStageCopy {
	does: string
	label: string
	leaves: string
}

export interface ScenarioCopy extends Record<LifecyclePhaseId, string> {
	context: string
	tabLabel: string
}

export interface OwnershipRowCopy {
	concern: string
	consumer: string
	platform: string
}

export interface StateCardCopy extends TitledCopy {
	eyebrow: string
}

export interface EvidenceStageCopy {
	kind: string
	title: string
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
			run: Record<RunMessageId, string>
		}
		card: {
			readMark: string
			notSet: string
			resolving: string
			labels: Record<PlanPhaseId, string>
			ticket: {
				status: string
				title: string
				links: string
				artifacts: Record<DocumentId, string>
			}
			analysis: {
				documents: string
				extracted: Record<DocumentId, string>
				guardrailJudgment: string
				documentsSummary: string
				episodes: string
				threshold: string
				verdicts: Record<EpisodeVerdict, string>
				reasons: Record<DropReason, string>
				episodesSummary: string
				measured: string
				measuredNone: string
				measuredReasoning: string
				suggestedModel: string
				suggestedEffort: string
				efforts: Record<PlanEffort, string>
				recommendationNote: string
				suggestionSummary: string
			}
			plan: {
				proposed: string
				pendingAuthorization: string
				envelope: { cpu: string; memory: string }
				ofBounded: string
				components: string
				columns: { cpu: string; memory: string }
				kinds: Record<ComponentKind, string>
				resolvedProgress: string
				componentsTotal: string
				crossCheck: string
				report: Record<
					| 'bounded'
					| 'headroom'
					| 'suggestedResource'
					| 'limits'
					| 'price'
					| 'validation'
					| 'evidence',
					string
				>
				deployments: string
				notRuntimeFields: string
				priceBasis: string
				validated: string
				planPath: string
				summary: string
			}
			workflow: {
				humanInTheLoop: string
				stages: Record<WorkflowStageId, string>
				humanStage: string
				mergeLocked: string
				requiredStage: string
				waiting: string
				personDriven: string
				personWriting: string
				model: string
				effort: string
				modelPrice: string
				estimate: string
				tokens: string
				reserves: string
				overBudget: string
				play: string
				starting: string
				playing: string
				playAgain: string
				playUnavailable: string
				result: Record<'reserved' | 'observed' | 'released', string>
				nothingPaid: string
				pullRequest: string
				safeSeat: string
				safeSeatMet: string
				checksRunning: string
				checksPassed: string
				readyForReview: string
				readyToMerge: string
				needsDecision: string
				summaryChoice: string
				summaryPaused: string
				summaryDone: string
				statusDone: string
				statusPaused: string
			}
			control: { pause: string; resume: string; replay: string }
			figures: string
			cpuValue: string
			memoryValue: string
		}
	}
	lifecycle: {
		stepperLabel: string
		stages: Record<LifecycleStageId, LifecycleStageCopy>
		whatItDoes: string
		evidenceItLeaves: string
	}
	scenarios: {
		title: string
		lead: string
		selectorLabel: string
		phases: Record<LifecyclePhaseId, string>
		items: Record<LifecycleScenarioId, ScenarioCopy>
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
		rows: Record<OwnershipRowId, OwnershipRowCopy>
	}
	state: {
		kicker: string
		title: string
		lead: string
		cards: Record<StateCardId, StateCardCopy>
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
		stages: Record<LifecycleStageId, EvidenceStageCopy>
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
		items: Record<MicroruntimesFaqId, FaqItemCopy>
	}
}
