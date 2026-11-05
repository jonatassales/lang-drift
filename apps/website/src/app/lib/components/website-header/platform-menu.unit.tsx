// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'

vi.mock('@i18n/navigation', () => ({
	Link: (props: {
		children: ReactNode
		className?: string
		href: string
		onBlur?: () => void
		onClick?: () => void
	}) => (
		<a className={props.className} href={props.href}>
			{props.children}
		</a>
	),
}))

vi.mock('next-intl', () => ({
	useTranslations: () => (key: string) => {
		const messages: Record<string, string> = {
			platform: 'Platform',
			platformMenu: 'Platform products',
			finops: 'FinOps',
			finopsDescription:
				'Fund the how: agent-assisted delivery accountable to a team budget.',
			microruntimes: 'Microruntimes',
			microruntimesDescription:
				'Bounded, consumer-owned execution around the software you already have.',
			cli: 'LangDrift CLI',
			cliDescription:
				'Plan, run and record bounded delivery from your terminal.',
			ciGates: 'CI Gates',
			ciGatesDescription:
				'Pull request reports and delivery evidence in your own CI.',
			safeSeat: 'Safe Seat',
			safeSeatDescription:
				'Part of LangDrift FinOps: every ticket keeps a part a person writes by hand.',
			comingSoon: 'Coming soon',
		}

		return messages[key] ?? key
	},
}))

import { PlatformMenu } from './platform-menu'

afterEach(() => {
	cleanup()
})

describe('website/src/app/lib/components/website-header/platform-menu', () => {
	it('lists the engine first and links every product, Safe Seat alone coming soon', () => {
		render(<PlatformMenu />)

		const hrefs = screen
			.getAllByRole('link')
			.map((link) => link.getAttribute('href'))

		expect(hrefs).toEqual([
			'/platform/microruntimes',
			'/platform/finops',
			'/platform/cli',
			'/platform/ci-gates',
			'/platform/safe-seat',
		])
		expect(screen.getAllByText('Coming soon')).toHaveLength(1)
	})

	it('describes CI Gates and Safe Seat', () => {
		render(<PlatformMenu />)

		expect(
			screen.getByText(
				'Pull request reports and delivery evidence in your own CI.',
			),
		).toBeDefined()
		expect(
			screen.getByText(
				'Part of LangDrift FinOps: every ticket keeps a part a person writes by hand.',
			),
		).toBeDefined()
		expect(screen.queryByText('Voice AI Interfaces')).toBeNull()
	})

	it('links Safe Seat to its page while it is coming soon', () => {
		render(<PlatformMenu />)

		const link = screen.getByText('Safe Seat').closest('a')

		expect(link?.getAttribute('href')).toBe('/platform/safe-seat')
	})
})
