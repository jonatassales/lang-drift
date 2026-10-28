// @vitest-environment jsdom
import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import {
	conversionMessages,
	finopsMessages,
	microruntimesMessages,
} from '@i18n/messages/index'

import { FaqTabs } from './faq-tabs'

const vision = conversionMessages.en.faq.items
const microruntimes = microruntimesMessages.en.faq.items
const finops = finopsMessages.en.faq.items

const groups = [
	{
		id: 'vision',
		items: [
			{
				id: 'meaning',
				...vision.meaning,
			},
			{
				id: 'tools',
				...vision.tools,
			},
		],
		label: 'Vision',
	},
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

describe('the FAQ by product', () => {
	it('opens on the first product and shows the next one when its tab is selected', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		expect(
			screen.getByRole('button', { name: vision.meaning.question }),
		).toBeTruthy()
		expect(
			screen.queryByRole('button', { name: microruntimes.what.question }),
		).toBeNull()

		fireEvent.click(screen.getByRole('tab', { name: 'Microruntimes' }))

		expect(
			screen.getByRole('button', { name: microruntimes.what.question }),
		).toBeTruthy()
		expect(
			screen.queryByRole('button', { name: vision.meaning.question }),
		).toBeNull()
		expect(
			screen
				.getByRole('tab', { name: 'Microruntimes' })
				.getAttribute('aria-selected'),
		).toBe('true')
	})

	it('shows a product collapsed again after switching away and back', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		const question = () =>
			screen.getByRole('button', { name: vision.meaning.question })
		fireEvent.click(question())
		expect(question().getAttribute('aria-expanded')).toBe('true')
		expect(screen.getByText(vision.meaning.answer)).toBeTruthy()

		fireEvent.click(screen.getByRole('tab', { name: 'FinOps' }))
		fireEvent.click(screen.getByRole('tab', { name: 'Vision' }))

		expect(question().getAttribute('aria-expanded')).toBe('false')
	})

	it('moves the selection with the arrow keys', () => {
		render(<FaqTabs groups={groups} tabsLabel="Product" />)

		const visionTab = screen.getByRole('tab', { name: 'Vision' })
		visionTab.focus()
		fireEvent.keyDown(visionTab, { key: 'ArrowRight' })

		expect(
			screen
				.getByRole('tab', { name: 'Microruntimes' })
				.getAttribute('aria-selected'),
		).toBe('true')
		expect(
			screen.getByRole('button', { name: microruntimes.what.question }),
		).toBeTruthy()
	})
})
