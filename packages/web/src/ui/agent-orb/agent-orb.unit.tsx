import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { AgentOrb } from './agent-orb'

describe('AgentOrb public props', () => {
	it('renders an SSR-safe native host with caller configuration and reserved dimensions', () => {
		const html = renderToString(
			<AgentOrb
				paused
				preset="peach"
				reducedMotion="always"
				size={92}
				speed={1.15}
				state="listening"
			/>,
		)
		const document = new DOMParser().parseFromString(html, 'text/html')
		const orb = document.querySelector('orb-z')

		expect(orb?.getAttribute('preset')).toBe('peach')
		expect(orb?.hasAttribute('paused')).toBe(true)
		expect(orb?.getAttribute('reduced-motion')).toBe('always')
		expect(orb?.getAttribute('size')).toBe('92px')
		expect(orb?.getAttribute('speed')).toBe('1.15')
		expect(orb?.getAttribute('state')).toBe('listening')
		expect(orb?.parentElement?.style.width).toBe('92px')
		expect(orb?.parentElement?.style.height).toBe('92px')
	})

	it('keeps default motion enabled and normalizes invalid numeric dimensions', () => {
		const html = renderToString(<AgentOrb size={-1} />)
		const document = new DOMParser().parseFromString(html, 'text/html')
		const orb = document.querySelector('orb-z')

		expect(orb?.hasAttribute('paused')).toBe(false)
		expect(orb?.getAttribute('reduced-motion')).toBe('system')
		expect(orb?.getAttribute('size')).toBe('16rem')
		expect(orb?.getAttribute('state')).toBe('idle')
	})
})
