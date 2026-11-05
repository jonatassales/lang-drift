import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { TextLink } from './text-link'

describe('web/ui/text-link/text-link', () => {
	it('renders an underlined anchor in the color around it', () => {
		const html = renderToString(<TextLink href="#sources">Sources</TextLink>)

		expect(html).toMatch(/^<a/)
		expect(html).toContain('href="#sources"')
		expect(html).toContain('underline')
		expect(html).toContain('text-inherit')
	})

	it('lends its look to its one child', () => {
		const html = renderToString(
			<TextLink asChild className="text-xs">
				<button type="button">Pause</button>
			</TextLink>,
		)

		expect(html).toMatch(/^<button/)
		expect(html).toContain('text-xs')
		expect(html).not.toContain('text-[13px]')
	})
})
