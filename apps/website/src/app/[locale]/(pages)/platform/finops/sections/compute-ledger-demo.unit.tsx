// @vitest-environment jsdom

import { cleanup, fireEvent, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { mountComputeLedgerDemo } from '@routes/platform/finops/sections/testdata/mount-compute-ledger-demo.mock'
import { selectedEnvironmentAtom } from '@state'

afterEach(cleanup)

describe('website/platform/finops/sections/compute-ledger-demo', () => {
	it('selects an environment and shows its receipt', () => {
		const store = mountComputeLedgerDemo()
		const button = screen.getByRole('button', { name: /Preview cluster/ })

		fireEvent.click(button)

		expect(store.get(selectedEnvironmentAtom)).toBe('preview-cluster')
		expect(button.getAttribute('aria-pressed')).toBe('true')

		const receipt = screen
			.getByRole('heading', {
				level: 2,
				name: 'Preview cluster',
			})
			.closest('article')

		expect(receipt?.textContent).toContain('2,210 CPU·h')
		expect(receipt?.textContent).toContain('96 h · $640')
		expect(receipt?.textContent).toContain('Still running after its run')
	})

	it('lists the CPU, memory and GPU hours of each metered environment', () => {
		mountComputeLedgerDemo()

		const button = screen.getByRole('button', { name: /CI runners/ })

		expect(button.textContent).toContain('1,480 CPU·h · 2,960 GB·h · 0 GPU·h')
	})

	it('shows an unmetered environment as unknown, never as zero', () => {
		mountComputeLedgerDemo()

		const button = screen.getByRole('button', { name: /Embedding GPU job/ })

		fireEvent.click(button)

		const receipt = screen
			.getByRole('heading', {
				level: 2,
				name: 'Embedding GPU job',
			})
			.closest('article')

		expect(button.textContent).toContain('meter not reporting')
		expect(button.textContent).not.toContain('CPU·h')
		expect(receipt?.textContent).toContain('Unknown · meter not reporting')
		expect(receipt?.textContent).not.toMatch(/\$0\b|\b0 (CPU|GB|GPU)·h/)
	})
})
