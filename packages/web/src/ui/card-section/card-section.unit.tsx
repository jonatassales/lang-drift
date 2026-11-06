import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { CardSection } from './card-section'

describe('web/ui/card-section/card-section', () => {
	it('sets its children under a hairline rule', () => {
		const html = renderToString(<CardSection>Episodes</CardSection>)

		expect(html).toBe(
			'<div class="mt-5 border-t border-border pt-4">Episodes</div>',
		)
	})

	it('lets the caller replace its spacing', () => {
		const html = renderToString(
			<CardSection className="mt-4 pt-3">Episodes</CardSection>,
		)

		expect(html).toContain('class="border-t border-border mt-4 pt-3"')
	})
})
