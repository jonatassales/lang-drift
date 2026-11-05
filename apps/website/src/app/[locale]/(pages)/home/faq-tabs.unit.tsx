// @vitest-environment jsdom
import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { finopsMessages, microruntimesMessages } from '@i18n/messages/index'

import { FaqTabs } from './faq-tabs'

const microruntimes = microruntimesMessages.en.faq.items
const finops = finopsMessages.en.faq.items

const groups = [
	{
		id: 'microruntimes',
		items: [
			{
				id: 'what',
				...microruntimes.what,
			},
			{
				id: 'sandbox',
				...microruntimes.sandbox,
			},
		],
		label: 'Microruntimes',
	},
	{
		id: 'finops',
		items: [
			{
				id: 'budget',
				...finops.budget,
			},
		],
		label: 'FinOps',
	},
]

afterEach(() => {
	cleanup()
})

describe('website/home/faq-tabs', () => {
	it('opens on the first product and shows the next one when its tab is selected', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		expect(
			screen.getByRole('button', { name: microruntimes.what.question }),
		).toBeTruthy()
		expect(
			screen.queryByRole('button', { name: finops.budget.question }),
		).toBeNull()

		fireEvent.click(screen.getByRole('tab', { name: 'FinOps' }))

		expect(
			screen.getByRole('button', { name: finops.budget.question }),
		).toBeTruthy()
		expect(
			screen.queryByRole('button', { name: microruntimes.what.question }),
		).toBeNull()
		expect(
			screen.getByRole('tab', { name: 'FinOps' }).getAttribute('aria-selected'),
		).toBe('true')
	})

	it('shows a product collapsed again after switching away and back', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		const question = () =>
			screen.getByRole('button', { name: microruntimes.what.question })
		fireEvent.click(question())
		expect(question().getAttribute('aria-expanded')).toBe('true')
		expect(screen.getByText(microruntimes.what.answer)).toBeTruthy()

		fireEvent.click(screen.getByRole('tab', { name: 'FinOps' }))
		fireEvent.click(screen.getByRole('tab', { name: 'Microruntimes' }))

		expect(question().getAttribute('aria-expanded')).toBe('false')
	})

	it('moves the selection with the arrow keys', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		const microruntimesTab = screen.getByRole('tab', { name: 'Microruntimes' })
		microruntimesTab.focus()
		fireEvent.keyDown(microruntimesTab, { key: 'ArrowRight' })

		expect(
			screen.getByRole('tab', { name: 'FinOps' }).getAttribute('aria-selected'),
		).toBe('true')
		expect(
			screen.getByRole('button', { name: finops.budget.question }),
		).toBeTruthy()
	})
})
