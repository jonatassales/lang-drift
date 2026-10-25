import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { SectionHeading } from './section-heading'

describe('SectionHeading', () => {
	it('renders the kicker before a heading that carries the id', () => {
		const html = renderToString(
			<SectionHeading
				headingId="budget-heading"
				kicker="Know who funds it"
				lead="Work varies in complexity."
				title="The budget belongs to the team."
			/>,
		)
		expect(html).toContain('id="budget-heading"')
		expect(html.indexOf('Know who funds it')).toBeLessThan(
			html.indexOf('The budget belongs to the team.'),
		)
		expect(html).toContain('Work varies in complexity.')
	})

	it('leaves the lead out when none is given', () => {
		const html = renderToString(
			<SectionHeading headingId="h" kicker="Kicker" title="Title" />,
		)
		expect(html).not.toContain('<p')
	})
})
