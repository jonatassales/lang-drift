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
import { headerMessages } from '@i18n/messages/index'

import { HomeProductSection } from './home-product-section'

const vision = platformProduct('vision')
const finops = platformProduct('finops')

afterEach(() => {
	cleanup()
})

describe('a home product section', () => {
	it('names the product, is labelled by its heading and ends with the link', () => {
		render(
			<HomeProductSection
				description="LangDrift shows how your product moved."
				headline="Why did your Product Vision fall?"
				id="vision"
				link="Explore Vision"
				locale="en"
				product={vision}
			>
				<p>The demo</p>
			</HomeProductSection>,
		)

		const region = screen.getByRole('region', {
			name: 'Why did your Product Vision fall?',
		})
		const demo = screen.getByText('The demo')
		const link = screen.getByRole('link', { name: 'Explore Vision' })

		expect(region.textContent).toContain(headerMessages.en.vision)
		expect(
			screen.getByRole('heading', {
				level: 2,
				name: 'Why did your Product Vision fall?',
			}).id,
		).toBe('vision-heading')
		expect(
			screen.getByText('LangDrift shows how your product moved.'),
		).toBeTruthy()
		expect(link.getAttribute('href')).toBe('/platform/vision')
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

		expect(region.textContent).toContain(headerMessages['pt-BR'].finops)
		expect(describedBy).toEqual([
			'Illustrative ledger.',
			'An illustrative team budget.',
		])
		expect(screen.queryByRole('link')).toBeNull()
	})

	it('relates nothing it is not given', () => {
		render(
			<HomeProductSection
				description="Ask about a change."
				headline="Keep its context in view."
				id="ai-voice"
				locale="en"
				product={platformProduct('aiVoice')}
			>
				<p>The orb</p>
			</HomeProductSection>,
		)

		const region = screen.getByRole('region', {
			name: 'Keep its context in view.',
		})

		expect(region.getAttribute('aria-describedby')).toBeNull()
		expect(screen.queryByRole('link')).toBeNull()
	})
})
