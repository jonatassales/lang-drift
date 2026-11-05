import { requiredHumanStages } from './microruntimes-workflow.data'
import type { HumanStageId } from './microruntimes-workflow.types'

export function isRequiredHumanStage(stage: HumanStageId): boolean {
	return requiredHumanStages.includes(stage)
}
