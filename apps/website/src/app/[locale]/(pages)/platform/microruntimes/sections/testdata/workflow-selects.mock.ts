import { screen } from '@testing-library/react'

// The Workflow card's model and effort selects.
export function workflowSelects() {
	return screen.getAllByRole('combobox', { hidden: true })
}
