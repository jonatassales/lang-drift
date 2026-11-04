import { describe, expect, it } from 'vitest'

import { shouldDeadLetter } from './should-dead-letter.compute.js'

describe('events/transport/should-dead-letter.compute', () => {
	it('keeps handling an entry delivered up to the maximum', () => {
		expect(shouldDeadLetter(1, 5)).toBe(false)
		expect(shouldDeadLetter(5, 5)).toBe(false)
	})

	it('parks an entry delivered more times than the maximum', () => {
		expect(shouldDeadLetter(6, 5)).toBe(true)
	})
})
