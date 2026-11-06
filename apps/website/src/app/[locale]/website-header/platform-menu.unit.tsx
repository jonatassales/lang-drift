// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'

interface LinkMockProps {
	children: ReactNode
	className?: string
	href: string
	onBlur?: () => void
	onClick?: () => void
}

vi.mock('@i18n/navigation', () => ({
	Link: (props: LinkMockProps) => (
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
			llmops: 'LLMOps',
			llmopsDescription:
				'Fund the reasoning: model calls accountable to a team budget.',
			finops: 'FinOps',
			finopsDescription:
				'Meter the hardware: CPU, memory and GPU held to a team budget while they run.',
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
				'Part of LangDrift LLMOps: every ticket keeps a part a person writes by hand.',
			comingSoon: 'Coming soon',
		}

		return messages[key] ?? key
	},
}))

import { PlatformMenu } from './platform-menu'

afterEach(() => {
	cleanup()
})

describe('website/website-header/platform-menu', () => {
	it('lists the products before the utilities and links each, Safe Seat alone coming soon', () => {
		render(<PlatformMenu />)

		const hrefs = screen
			.getAllByRole('link')
			.map((link) => link.getAttribute('href'))

		expect(hrefs).toEqual([
			'/platform/llmops',
			'/platform/finops',
			'/platform/safe-seat',
			'/platform/microruntimes',
			'/platform/cli',
			'/platform/ci-gates',
		])
		expect(screen.getAllByText('Coming soon')).toHaveLength(1)
	})

	it('describes LLMOps, FinOps, CI Gates and Safe Seat', () => {
		render(<PlatformMenu />)

		expect(
			screen.getByText(
				'Fund the reasoning: model calls accountable to a team budget.',
			),
		).toBeDefined()
		expect(
			screen.getByText(
				'Meter the hardware: CPU, memory and GPU held to a team budget while they run.',
			),
		).toBeDefined()
		expect(
			screen.getByText(
				'Pull request reports and delivery evidence in your own CI.',
			),
		).toBeDefined()
		expect(
			screen.getByText(
				'Part of LangDrift LLMOps: every ticket keeps a part a person writes by hand.',
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
