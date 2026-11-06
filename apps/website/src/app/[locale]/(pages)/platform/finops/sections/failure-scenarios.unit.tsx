// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { finopsMessages } from '@i18n/messages/index'

import { FailureScenarios } from './failure-scenarios'

afterEach(cleanup)

describe('website/platform/finops/sections/failure-scenarios', () => {
	it('walks each hardware failure through the chain from intent to outcome', () => {
		const copy = finopsMessages.en.failures

		render(<FailureScenarios copy={copy} />)

		expect(
			screen
				.getAllByRole('heading', { level: 3 })
				.map((heading) => heading.textContent),
		).toEqual(['1Intent', '2Delivery', '3Evidence', '4Outcome'])
		expect(
			screen.getByText(copy.scenarios['orphaned-environment'].context),
		).toBeTruthy()

		fireEvent.click(screen.getByRole('tab', { name: 'Meter stopped' }))

		expect(
			screen.getByText(copy.scenarios['silent-meter'].outcome),
		).toBeTruthy()
	})
})
