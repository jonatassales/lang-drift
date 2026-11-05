// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { ProductBadge } from './product-badge'

afterEach(() => {
	cleanup()
})

describe('website/lib/components/product-badge', () => {
	it('renders children with discreet badge styling', () => {
		render(<ProductBadge>Coming soon</ProductBadge>)

		const badge = screen.getByText('Coming soon')

		expect(badge).toBeDefined()
		expect(badge.className).toContain('rounded-full')
		expect(badge.className).toContain('border')
		expect(badge.className).toContain('text-muted-foreground')
	})

	it('merges additional class names', () => {
		render(<ProductBadge className="custom-class">Coming soon</ProductBadge>)

		const badge = screen.getByText('Coming soon')

		expect(badge.className).toContain('custom-class')
	})
})
