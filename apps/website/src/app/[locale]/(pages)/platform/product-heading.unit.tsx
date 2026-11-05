// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { FinopsIcon } from 'web/ui/icons/platform'

import { ProductHeading } from './product-heading'

afterEach(() => {
	cleanup()
})

describe('website/platform/product-heading', () => {
	it('names the product in an h1 opened by its icon tile', () => {
		render(
			<ProductHeading icon={FinopsIcon} id="finops-hero-heading">
				FinOps
			</ProductHeading>,
		)

		const heading = screen.getByRole('heading', {
			level: 1,
			name: 'FinOps',
		})
		const tile = heading.firstElementChild
		const icon = tile?.firstElementChild

		expect(heading.id).toBe('finops-hero-heading')
		expect(tile?.getAttribute('aria-hidden')).toBe('true')
		expect(tile?.getAttribute('class')).toContain('size-[0.74em]')
		expect(icon?.getAttribute('class')).toContain('size-[0.46em]')
	})
})
