import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BasicAccordion } from './basic-accordion'

const items = [
	{
		content: 'A team envelope.',
		id: 'budget',
		title: 'What is the budget?',
	},
	{
		content: 'Unknown, never zero.',
		id: 'unknown',
		title: 'What if spend is unknown?',
	},
]

describe('web/vendors/smoothui/basic-accordion', () => {
	it('divides the items and rotates the open item chevron', () => {
		const html = renderToString(
			<BasicAccordion defaultExpandedIds={['budget']} items={items} />,
		)

		expect(html.match(/not-first:border-t/g)).toHaveLength(2)
		expect(html).toContain('aria-expanded="true"')
		expect(html).toContain('data-expanded="true"')
		expect(html).toContain('data-[expanded=true]:rotate-180')
		expect(html).toContain('motion-reduce:transition-none')
		expect(html).toContain('focus-visible:outline-foreground')
		expect(html).not.toContain('smooth-')
	})

	it('hides a closed item region', () => {
		const html = renderToString(<BasicAccordion items={items} />)

		expect(html.match(/hidden=""/g)).toHaveLength(2)
	})
})
