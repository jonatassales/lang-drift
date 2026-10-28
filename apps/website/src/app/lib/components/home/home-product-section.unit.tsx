// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'

vi.mock('@i18n/navigation', () => ({
	Link: (props: { children: ReactNode; href: string }) => (
		<a href={props.href}>{props.children}</a>
	),
}))

import { platformProduct } from '@components/website-header/platform-products'

import { HomeProductSection } from './home-product-section'

const microruntimes = platformProduct('microruntimes')
const finops = platformProduct('finops')

afterEach(() => {
	cleanup()
})

describe('a home product section', () => {
	it('is labelled by its heading and ends with the link', () => {
		render(
			<HomeProductSection
				description="A run is bounded before it starts."
				headline="How does a run stay bounded?"
				id="microruntimes"
				link="Explore Microruntimes"
				locale="en"
				product={microruntimes}
			>
				<p>The demo</p>
			</HomeProductSection>,
		)

		const demo = screen.getByText('The demo')
		const link = screen.getByRole('link', { name: 'Explore Microruntimes' })

		expect(
			screen.getByRole('region', {
				name: 'How does a run stay bounded?',
			}),
		).toBeTruthy()
		expect(
			screen.getByRole('heading', {
				level: 2,
				name: 'How does a run stay bounded?',
			}).id,
		).toBe('microruntimes-heading')
		expect(screen.getByText('A run is bounded before it starts.')).toBeTruthy()
		expect(link.getAttribute('href')).toBe('/platform/microruntimes')
		expect(
			demo.compareDocumentPosition(link) & Node.DOCUMENT_POSITION_FOLLOWING,
		).toBeTruthy()
	})

	it('is described by its caption and its summary when the section gives them', () => {
		render(
			<HomeProductSection
				caption="Illustrative ledger."
				description="Every paid call is reserved before it runs."
				headline="One team envelope."
				id="finops"
				locale="pt-BR"
				product={finops}
				summary="An illustrative team budget."
			>
				<p>The ledger</p>
			</HomeProductSection>,
		)

		const region = screen.getByRole('region', { name: 'One team envelope.' })
		const describedBy = (region.getAttribute('aria-describedby') ?? '')
			.split(' ')
			.map((id) => document.getElementById(id)?.textContent)

		expect(describedBy).toEqual([
			'Illustrative ledger.',
			'An illustrative team budget.',
		])
		expect(screen.queryByRole('link')).toBeNull()
	})

	it('relates nothing it is not given', () => {
		render(
			<HomeProductSection
				description="A run stays inside its envelope."
				headline="Bounded before it starts."
				id="microruntimes"
				locale="en"
				product={microruntimes}
			>
				<p>The session</p>
			</HomeProductSection>,
		)

		const region = screen.getByRole('region', {
			name: 'Bounded before it starts.',
		})

		expect(region.getAttribute('aria-describedby')).toBeNull()
		expect(screen.queryByRole('link')).toBeNull()
	})
})
