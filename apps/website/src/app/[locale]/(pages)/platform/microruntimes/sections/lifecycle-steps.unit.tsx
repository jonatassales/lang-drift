// @vitest-environment jsdom

import { cleanup, fireEvent, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { mountLifecycleSteps } from '@routes/platform/microruntimes/sections/testdata/mount-lifecycle-steps.mock'
import { selectedLifecycleStageAtom } from '@state'

afterEach(cleanup)

describe('website/platform/microruntimes/sections/lifecycle-steps', () => {
	it('moves the current step and the evidence record together', () => {
		const store = mountLifecycleSteps()
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
