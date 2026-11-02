import {
	type WorkflowModel,
	type WorkflowModelId,
	workflowModels,
} from './microruntimes-workflow'

export function workflowModel(id: WorkflowModelId): WorkflowModel {
	const model = workflowModels.find((candidate) => candidate.id === id)
	if (!model) throw new Error(`Unknown workflow model: ${id}`)

	return model
}
