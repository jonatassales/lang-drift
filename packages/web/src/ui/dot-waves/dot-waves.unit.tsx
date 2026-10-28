import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { DotWaves } from './dot-waves'

describe('DotWaves markup', () => {
	it('renders a hidden, pointer-transparent layer with one canvas that fills its parent and merges classes', () => {
		const html = renderToString(<DotWaves className="opacity-80" />)
		const document = new DOMParser().parseFromString(html, 'text/html')
		const layer = document.body.firstElementChild

		expect(layer?.getAttribute('aria-hidden')).toBe('true')
		expect(layer?.className).toContain('pointer-events-none')
		expect(layer?.className).toContain('absolute')
		expect(layer?.className).toContain('opacity-80')
		expect(layer?.querySelectorAll('canvas')).toHaveLength(1)
	})
})
