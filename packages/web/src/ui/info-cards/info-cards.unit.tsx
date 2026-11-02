import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { InfoCards } from './info-cards'

const items = [
	{
		description: 'One authorized amount.',
		eyebrow: 'Team envelope',
		id: 'envelope',
		title: 'One amount, one period',
	},
	{
		description: 'A suggested share.',
		id: 'allocation',
		title: 'Allocation',
	},
]

describe('web/ui/info-cards/info-cards', () => {
	it('renders an unordered list of cards with optional eyebrows', () => {
		const html = renderToString(<InfoCards items={items} />)
		expect(html.startsWith('<ul')).toBe(true)
		expect(html).toContain('Team envelope')
		expect(html).toContain('lg:grid-cols-3')
		expect(html.match(/<h3/g)).toHaveLength(2)
	})

	it('renders steps as an ordered list with the requested columns', () => {
		const html = renderToString(<InfoCards columns={4} items={items} ordered />)
		expect(html.startsWith('<ol')).toBe(true)
		expect(html).toContain('lg:grid-cols-4')
	})

	it('renders the footer the caller passes', () => {
		const html = renderToString(
			<InfoCards
				items={[
					{
						...items[1],
						footer: <a href="/platform/microruntimes">See it</a>,
					},
				]}
			/>,
		)
		expect(html).toContain('href="/platform/microruntimes"')
	})
})
