// @vitest-environment jsdom

import { cleanup, render, screen, within } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { llmopsMessages } from '@i18n/messages/index'

import { FailureScenarios } from './failure-scenarios'

afterEach(cleanup)

describe('website/platform/llmops/sections/failure-scenarios', () => {
	it('walks both failures through the chain and ends each in the state it leaves', () => {
		const { container } = render(
			<FailureScenarios copy={llmopsMessages.en.failures} />,
		)
		const [pause, unknown] = screen.getAllByRole('article')

		expect(within(pause).getByRole('heading', { level: 3 }).textContent).toBe(
			'Pause and replan',
		)
		expect(
			within(pause)
				.getAllByRole('term')
				.map((term) => term.textContent),
		).toEqual(['Intent', 'Delivery', 'Evidence', 'Outcome', 'Needs attention'])
		expect(within(pause).getAllByRole('definition').at(-1)?.textContent).toBe(
			'The reindexing job is paused before its next retry would exceed the reservation; a decision is pending.',
		)

		expect(within(unknown).getByRole('heading', { level: 3 }).textContent).toBe(
			'Unknown after a timeout',
		)
		expect(within(unknown).getAllByRole('definition').at(-1)?.textContent).toBe(
			'Unknown · awaiting reconciliation',
		)
		expect(container.textContent).not.toMatch(/\$0\b/)
	})
})
