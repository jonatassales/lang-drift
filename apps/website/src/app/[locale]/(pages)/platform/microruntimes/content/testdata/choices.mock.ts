import { initialWorkflowChoices } from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import type { WorkflowChoices } from '@routes/platform/microruntimes/content/microruntimes-workflow.types'

export function choices(overrides: Partial<WorkflowChoices>): WorkflowChoices {
	return {
		...initialWorkflowChoices,
		...overrides,
	}
}
