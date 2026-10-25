// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, describe, expect, it } from 'vitest'

import { microruntimesMessages } from '@i18n/messages/index'
import { selectedLifecycleStageAtom } from '@state'

import { LifecycleSteps } from './lifecycle-steps'
import { StageEvidence } from './stage-evidence'

afterEach(cleanup)

function mount() {
	const store = createStore()
	const copy = microruntimesMessages.en
	render(
		<Provider store={store}>
			<LifecycleSteps copy={copy.lifecycle} />
			<StageEvidence copy={copy.evidence} stageLabels={copy.lifecycle.stages} />
		</Provider>,
	)

	return store
}

describe('Microruntimes lifecycle steps', () => {
	it('moves the current step and the evidence record together', () => {
		const store = mount()
		const prepare = screen.getByRole('button', { name: /Prepare/ })
		expect(prepare.getAttribute('aria-current')).toBeNull()
		fireEvent.click(prepare)
		expect(store.get(selectedLifecycleStageAtom)).toBe('prepare')
		expect(prepare.getAttribute('aria-current')).toBe('step')
		expect(
			screen.getByText('Namespace and quota recorded before creation'),
		).toBeTruthy()
	})
})
