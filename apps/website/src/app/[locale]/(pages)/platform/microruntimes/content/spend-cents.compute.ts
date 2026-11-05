import type { Tokens, WorkflowModelId } from './microruntimes-workflow.types'
import { workflowModel } from './workflow-model.compute'

// What the tokens cost at the model's uncached list prices.
export function spendCents(modelId: WorkflowModelId, tokens: Tokens): number {
	const model = workflowModel(modelId)
	const tenthsOfCents =
		tokens.inputThousands * model.inputCents +
		tokens.outputThousands * model.outputCents

	return Math.round(tenthsOfCents / 1000)
}
