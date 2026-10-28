// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { createStore, Provider } from 'jotai'
import { afterEach, describe, expect, it } from 'vitest'

import { finopsMessages } from '@i18n/messages/index'
import { selectedWorkItemAtom } from '@state'

import { BudgetLedgerDemo } from './budget-ledger-demo'
import { DeliveryReceipt } from './delivery-receipt'

afterEach(cleanup)

function mount() {
	const store = createStore()
	const copy = finopsMessages.en
	render(
		<Provider store={store}>
			<BudgetLedgerDemo copy={copy.ledger} locale="en" />
			<DeliveryReceipt
				copy={copy.receipt}
				items={copy.ledger.items}
				locale="en"
			/>
		</Provider>,
	)

	return store
}

describe('FinOps ledger and receipt', () => {
	it('selects a work item and shows its receipt', () => {
		const store = mount()
		const button = screen.getByRole('button', { name: /Tenant migration/ })
		fireEvent.click(button)
		expect(store.get(selectedWorkItemAtom)).toBe('tenant-migration')
		expect(button.getAttribute('aria-pressed')).toBe('true')
		expect(
			screen.getByRole('heading', {
				level: 2,
				name: 'Tenant migration',
			}),
		).toBeTruthy()
	})

	it('shows unknown spend as unknown, never as zero', () => {
		mount()
		fireEvent.click(screen.getByRole('button', { name: /Invoice export job/ }))
		const receipt = screen
			.getByRole('heading', {
				level: 2,
				name: 'Invoice export job',
			})
			.closest('article')
		expect(receipt?.textContent).toContain('Unknown · awaiting reconciliation')
		expect(receipt?.textContent).not.toMatch(/\$0\b/)
	})

	it('formats the model tier and effort through the catalog', () => {
		mount()
		expect(screen.getByText('Deep reasoning · High')).toBeTruthy()
	})

	it('counts retries with the plural copy', () => {
		mount()
		fireEvent.click(screen.getByRole('button', { name: /Onboarding copy/ }))
		expect(screen.getByText('No retries')).toBeTruthy()
		fireEvent.click(screen.getByRole('button', { name: /Checkout retries/ }))
		expect(screen.getByText('2 retries')).toBeTruthy()
	})
})
