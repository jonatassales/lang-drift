// @vitest-environment jsdom

import { cleanup, fireEvent, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { mountBudgetLedgerDemo } from '@routes/platform/finops/sections/testdata/mount-budget-ledger-demo.mock'
import { selectedWorkItemAtom } from '@state'

afterEach(cleanup)

describe('website/platform/finops/sections/budget-ledger-demo', () => {
	it('selects a work item and shows its receipt', () => {
		const store = mountBudgetLedgerDemo()
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
		mountBudgetLedgerDemo()
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
		mountBudgetLedgerDemo()
		expect(screen.getByText('Deep reasoning · High')).toBeTruthy()
	})

	it('counts retries with the plural copy', () => {
		mountBudgetLedgerDemo()
		fireEvent.click(screen.getByRole('button', { name: /Onboarding copy/ }))
		expect(screen.getByText('No retries')).toBeTruthy()
		fireEvent.click(screen.getByRole('button', { name: /Checkout retries/ }))
		expect(screen.getByText('2 retries')).toBeTruthy()
	})
})
