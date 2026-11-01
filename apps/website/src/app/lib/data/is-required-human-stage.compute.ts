import {
	type HumanStageId,
	requiredHumanStages,
} from './microruntimes-workflow'

export function isRequiredHumanStage(stage: HumanStageId): boolean {
	return requiredHumanStages.includes(stage)
}
