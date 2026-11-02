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
			ciGates: 'CI Gates',
			ciGatesDescription:
				'Pull request reports and delivery evidence in your own CI.',
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
	it('renders platform products with CI Gates coming soon badge and description', () => {
		render(<PlatformMenu />)

		expect(screen.getByText('CI Gates')).toBeDefined()
		expect(screen.getByText('Coming soon')).toBeDefined()
		expect(
			screen.getByText(
				'Pull request reports and delivery evidence in your own CI.',
			),
		).toBeDefined()
	})
})
