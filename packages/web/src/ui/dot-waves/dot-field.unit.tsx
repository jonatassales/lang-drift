import { describe, expect, it } from 'vitest'

import { type Dot, dashOf, layoutDots, spacing, waveReach } from './dot-field'

const center = {
	left: 400,
	top: 300,
}

function dotAt(left: number, top: number): Dot {
	return {
		left,
		phase: 0,
		top,
	}
}

describe('the dot field', () => {
	it('lays one dot per cell, jittered inside it, the same way every time', () => {
		const dots = layoutDots(spacing * 4, spacing * 3)
		const again = layoutDots(spacing * 4, spacing * 3)

		expect(dots).toHaveLength(12)
		expect(again).toEqual(dots)
		for (const dot of dots) {
			expect(dot.left).toBeGreaterThanOrEqual(0)
			expect(dot.left).toBeLessThanOrEqual(spacing * 4)
			expect(dot.top).toBeGreaterThanOrEqual(0)
			expect(dot.top).toBeLessThanOrEqual(spacing * 3)
		}
	})

	it('points every dash away from the wave center', () => {
		const right = dashOf(dotAt(500, 300), center, 0)
		const above = dashOf(dotAt(400, 200), center, 0)

		expect(right.toLeft).toBeGreaterThan(0)
		expect(Math.abs(right.toTop)).toBeLessThan(1e-9)
		expect(above.toTop).toBeLessThan(0)
		expect(Math.abs(above.toLeft)).toBeLessThan(1e-9)
	})

	it('moves the rings outward as time passes', () => {
		const dot = dotAt(500, 300)
		const strengths = [0, 0.25, 0.5, 0.75, 1].map(
			(time) => dashOf(dot, center, time).alpha,
		)

		expect(new Set(strengths).size).toBeGreaterThan(1)
	})

	it('keeps a faint dot beyond the reach of the waves', () => {
		const far = dashOf(dotAt(400 + waveReach + 50, 300), center, 0)

		expect(far.alpha).toBeGreaterThan(0)
		expect(far.alpha).toBeLessThan(0.25)
	})
})
