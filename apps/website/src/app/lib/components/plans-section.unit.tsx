// @vitest-environment jsdom
import { cleanup, render, screen, within } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'

interface LinkProps {
	children: ReactNode
	className?: string
	href: string
}

vi.mock('@i18n/navigation', () => ({
	Link: (props: LinkProps) => <a href={props.href}>{props.children}</a>,
}))

import { conversionMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import { PlansSection } from './plans-section'

afterEach(() => {
	cleanup()
})

describe('website/lib/components/plans-section', () => {
	it.each(routing.locales)(
		'shows Starter, Plus and Pro, to be announced, in %s',
		(locale) => {
			const copy = conversionMessages[locale].plans

			render(<PlansSection locale={locale} />)

			const cards = screen.getAllByRole('article')

			expect(cards.map((card) => card.dataset.plan)).toEqual([
				'starter',
				'plus',
				'pro',
			])

			for (const card of cards) {
				const plan = copy.cards[card.dataset.plan as 'starter' | 'plus' | 'pro']

				expect(
					within(card).getByRole('heading', { level: 2 }).textContent,
				).toBe(plan.name)
				expect(within(card).getByText(plan.value)).toBeTruthy()
				expect(within(card).getByRole('link').getAttribute('href')).toBe(
					'/book-a-demo',
				)
			}

			expect(screen.getByText(copy.comingSoon)).toBeTruthy()
		},
	)
})
