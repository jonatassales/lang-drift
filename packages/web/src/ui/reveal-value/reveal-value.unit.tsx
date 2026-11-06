import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { RevealValue } from './reveal-value'

const labels = {
	notSet: 'Not set yet',
	resolving: 'Resolving',
}

describe('web/ui/reveal-value/reveal-value', () => {
	it('shows a value not known yet as a dash and says so', () => {
		const html = renderToString(
			<RevealValue className="font-mono" labels={labels} state="unset">
				250 m
			</RevealValue>,
		)

		expect(html).toBe(
			'<span class="font-mono"><span aria-hidden="true">—</span><span class="sr-only">Not set yet</span></span>',
		)
	})

	it('pulses while the value resolves and says so', () => {
		const html = renderToString(
			<RevealValue labels={labels} state="resolving">
				250 m
			</RevealValue>,
		)

		expect(html).toContain('animate-pulse motion-reduce:animate-none w-24')
		expect(html).toContain('<span class="sr-only">Resolving</span>')
		expect(html).not.toContain('250 m')
	})

	it('sizes the resolving bar as asked', () => {
		const html = renderToString(
			<RevealValue labels={labels} state="resolving" width="w-48" />,
		)

		expect(html).toContain('w-48')
		expect(html).not.toContain('w-24')
	})

	it('fades the value in once it is known', () => {
		const html = renderToString(
			<RevealValue className="font-mono" labels={labels} state="set">
				250 m
			</RevealValue>,
		)

		expect(html).toBe(
			'<span class="motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0 font-mono">250 m</span>',
		)
	})
})
