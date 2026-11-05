import {
	initialWorkflowChoices,
	type WorkflowChoices,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.data'

export function choices(overrides: Partial<WorkflowChoices>): WorkflowChoices {
	return {
		...initialWorkflowChoices,
		...overrides,
	}
}
