import { act, fireEvent } from '@testing-library/react'
import type { createStore } from 'jotai'

import { runScriptAtom, startRunAtom } from '@state'

import { workflowPlayButton } from './workflow-play-button.mock'

// Launches the run from the Workflow card and returns the script it plays.
export function launchWorkflowRun(store: ReturnType<typeof createStore>) {
	fireEvent.click(workflowPlayButton())
	act(() => {
		store.set(startRunAtom)
	})

	return store.get(runScriptAtom)
}
