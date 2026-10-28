import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { ElevatedFrame } from './elevated-frame'

describe('ElevatedFrame', () => {
	it('frames its children in glass and merges classes over its own', () => {
		const html = renderToString(
			<ElevatedFrame className="p-6 rounded-xl">
				<p>Inside</p>
			</ElevatedFrame>,
		)
		const frame = new DOMParser().parseFromString(html, 'text/html').body
			.firstElementChild

		expect(frame?.textContent).toBe('Inside')
		expect(frame?.classList.contains('backdrop-blur-xl')).toBe(true)
		expect(frame?.classList.contains('border-border')).toBe(true)
		expect(frame?.classList.contains('p-6')).toBe(true)
		expect(frame?.classList.contains('rounded-xl')).toBe(true)
		expect(frame?.classList.contains('rounded-2xl')).toBe(false)
	})
})
