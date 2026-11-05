import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { AnimatedTabs } from './animated-tabs'

const tabs = [
	{
		id: 'finops',
		label: 'FinOps',
	},
	{
		id: 'microruntimes',
		label: 'Microruntimes',
	},
]

describe('web/vendors/smoothui/animated-tabs', () => {
	it('draws an underline list with a primary bar under the selected tab', () => {
		const html = renderToString(
			<AnimatedTabs ariaLabel="Products" defaultTab="finops" tabs={tabs} />,
		)

		expect(html).toContain('data-variant="underline"')
		expect(html).toContain('border-b')
		expect(html).toContain('h-0.5 bg-primary')
		expect(html.match(/h-0\.5 bg-primary/g)).toHaveLength(1)
		expect(html).toContain('aria-selected:text-foreground')
		expect(html).toContain('focus-visible:outline-foreground')
		expect(html).toContain('motion-reduce:transition-none')
		expect(html).not.toContain('smooth-')
	})

	it('fills the selected pill or segment tab with the card role', () => {
		const pill = renderToString(
			<AnimatedTabs ariaLabel="Products" tabs={tabs} variant="pill" />,
		)
		const segment = renderToString(
			<AnimatedTabs ariaLabel="Products" tabs={tabs} variant="segment" />,
		)

		expect(pill).toContain('rounded-full border bg-background')
		expect(pill).toContain(
			'inset-0 rounded-[inherit] border border-border bg-card',
		)
		expect(segment).toContain('gap-0 rounded-lg border bg-background')
		expect(segment).toContain('rounded-md')
		expect(segment).toContain(
			'inset-0 rounded-[inherit] border border-border bg-card',
		)
	})
})
