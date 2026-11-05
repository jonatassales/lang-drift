import { workflowModels } from './microruntimes-workflow.data'
import type {
	WorkflowModel,
	WorkflowModelId,
} from './microruntimes-workflow.types'

export function workflowModel(id: WorkflowModelId): WorkflowModel {
	const model = workflowModels.find((candidate) => candidate.id === id)

	if (!model) throw new Error(`Unknown workflow model: ${id}`)

	return model
}
