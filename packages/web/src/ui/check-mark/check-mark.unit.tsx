import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { CheckMark } from './check-mark'

describe('web/ui/check-mark/check-mark', () => {
	it('hides the mark from assistive technology and says its label', () => {
		const html = renderToString(<CheckMark label="Read" />)

		expect(html).toContain('<span aria-hidden="true">✓</span>')
		expect(html).toContain('<span class="sr-only">Read</span>')
	})

	it('draws the mark in the success role', () => {
		const html = renderToString(<CheckMark label="Read" />)

		expect(html).toContain('text-status-success')
	})
})
