import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Button } from './button'

describe('web/ui/button', () => {
	it('renders a filled button by default', () => {
		const html = renderToString(<Button>Save</Button>)

		expect(html).toMatch(/^<button/)
		expect(html).toContain('bg-primary')
		expect(html).toContain('data-variant="default"')
	})

	it('draws the outline and destructive variants with token roles', () => {
		const outline = renderToString(<Button variant="outline">Retry</Button>)
		const destructive = renderToString(
			<Button variant="destructive">Delete</Button>,
		)

		expect(outline).toContain('border-border')
		expect(destructive).toContain('bg-destructive')
	})

	it('lets the caller override a class', () => {
		const html = renderToString(<Button className="h-12">Open</Button>)

		expect(html).toContain('h-12')
		expect(html).not.toContain('h-9')
	})

	it('lends its look to its child', () => {
		const html = renderToString(
			<Button asChild>
				<a href="#plans">Plans</a>
			</Button>,
		)

		expect(html).toMatch(/^<a/)
		expect(html).toContain('href="#plans"')
		expect(html).toContain('bg-primary')
	})
})
