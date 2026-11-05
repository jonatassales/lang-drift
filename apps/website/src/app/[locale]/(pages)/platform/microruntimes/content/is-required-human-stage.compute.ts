import {
	type HumanStageId,
	requiredHumanStages,
} from './microruntimes-workflow.data'

export function isRequiredHumanStage(stage: HumanStageId): boolean {
	return requiredHumanStages.includes(stage)
}
