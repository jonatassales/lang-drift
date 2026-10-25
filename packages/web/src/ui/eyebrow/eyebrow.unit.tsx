import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Eyebrow } from './eyebrow'

describe('Eyebrow', () => {
	it('renders a span by default and a term on request', () => {
		expect(renderToString(<Eyebrow>Period</Eyebrow>)).toMatch(/^<span/)
		expect(renderToString(<Eyebrow as="dt">Estimate</Eyebrow>)).toMatch(/^<dt/)
	})

	it('lets the caller adjust the label', () => {
		const html = renderToString(
			<Eyebrow className="text-[10px]">Kicker</Eyebrow>,
		)
		expect(html).toContain('text-[10px]')
		expect(html).not.toContain('text-[9px]')
	})
})
