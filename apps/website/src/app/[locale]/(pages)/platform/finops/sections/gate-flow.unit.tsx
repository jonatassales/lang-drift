// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { finopsMessages } from '@i18n/messages/index'
import {
	gateExample,
	gateFigures,
} from '@routes/platform/finops/content/finops-content'

import { GateFlow } from './gate-flow'

afterEach(cleanup)

function mount() {
	const copy = finopsMessages.en

	return render(
		<GateFlow
			copy={copy.authorization}
			exampleTitle={copy.ledger.items[gateExample.id].title}
			figures={gateFigures(gateExample)}
			locale="en"
		/>,
	)
}

describe('FinOps gate flow', () => {
	it('walks the three steps in order with the example receipt figures', () => {
		mount()
		const steps = screen
			.getAllByRole('listitem')
			.map((step) => step.textContent ?? '')
		expect(steps).toHaveLength(3)
		expect(steps[0]).toContain('01 · The agent recommends')
		expect(steps[0]).toContain('$1,800')
		expect(steps[1]).toContain('02 · The gate reserves')
		expect(steps[1]).toContain('$2,600')
		expect(steps[2]).toContain('03 · The evidence closes')
		expect(steps[2]).toContain('$2,460')
		expect(steps[2]).toContain('$140')
	})

	it('names the work item the figures come from and never shows zero', () => {
		const { container } = mount()
		expect(
			screen.getByText(
				/Figures from the Checkout retries receipt at the top of the page/,
			),
		).toBeTruthy()
		expect(container.textContent).not.toMatch(/\$0\b/)
	})
})
