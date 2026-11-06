import { describe, expect, it } from 'vitest'

import { lightningBolt } from './lightning-bolt.compute'

describe('web/ui/short-circuit-text/lightning-bolt.compute', () => {
	it('leaps from the word edge in six segments', () => {
		const path = lightningBolt(() => 0.5, 1)

		expect(path).toBe(
			'M100 50.0 L83.3 50.0 L66.7 50.0 L50.0 50.0 L33.3 50.0 L16.7 50.0 L0.0 50.0',
		)
	})

	it('never climbs past the top edge of the box', () => {
		const path = lightningBolt(() => 0, 1)
		const heights = [...path.matchAll(/ (-?[\d.]+)/g)].map((match) =>
			Number(match[1]),
		)

		expect(Math.min(...heights)).toBe(2)
	})

	it('never drops past the bottom edge of the box', () => {
		const path = lightningBolt(() => 1, 1)
		const heights = [...path.matchAll(/ (-?[\d.]+)/g)].map((match) =>
			Number(match[1]),
		)

		expect(Math.max(...heights)).toBe(98)
	})
})
