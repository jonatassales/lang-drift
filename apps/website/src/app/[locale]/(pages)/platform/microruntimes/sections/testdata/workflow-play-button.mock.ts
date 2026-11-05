import { screen } from '@testing-library/react'

// The Workflow card's play button, whatever stage its label shows.
export function workflowPlayButton() {
	return screen.getByRole('button', {
		hidden: true,
		name: /workflow|Starting|Playing|Play again/,
	})
}
