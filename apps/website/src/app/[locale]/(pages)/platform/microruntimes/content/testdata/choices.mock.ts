import {
	initialWorkflowChoices,
	type WorkflowChoices,
} from '@routes/platform/microruntimes/content/microruntimes-workflow'

export function choices(overrides: Partial<WorkflowChoices>): WorkflowChoices {
	return {
		...initialWorkflowChoices,
		...overrides,
	}
}
