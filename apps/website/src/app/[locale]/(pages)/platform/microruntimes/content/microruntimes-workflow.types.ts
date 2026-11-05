import type { TerminalLineType } from 'web/vendors/smoothui/mac-terminal'

import type { PlanEffort } from './microruntimes-content.types'
import type {
	runMessageIds,
	workflowModels,
	workflowStageIds,
} from './microruntimes-workflow.data'

export type WorkflowModel = (typeof workflowModels)[number]

export type WorkflowModelId = WorkflowModel['id']

// Input and output tokens, in thousands.
export interface Tokens {
	inputThousands: number
	outputThousands: number
}

export type WorkflowStageId = (typeof workflowStageIds)[number]

// The stages a person can be asked to confirm. The dev ticket is the
// person's own work, never switched off, and merge is always the owner's.
export type HumanStageId = Exclude<WorkflowStageId, 'dev-ticket' | 'merge'>

export interface WorkflowChoices {
	effort: PlanEffort
	human: Record<HumanStageId, boolean>
	model: WorkflowModelId
}

export type RunMessageId = (typeof runMessageIds)[number]

export interface RunLine {
	delay: number
	id: string
	message: RunMessageId
	stage?: WorkflowStageId
	type: TerminalLineType
	values?: Record<string, string>
}

export type StageProgress = 'pending' | 'active' | 'done'
