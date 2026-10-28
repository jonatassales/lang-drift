// The geometry of the dot waves, free of any canvas: where the dots sit, and
// how each one is drawn for a wave center and a moment in time.

export interface Point {
	left: number
	top: number
}

export interface Dot extends Point {
	/** A fixed offset of the wave for this dot, between 0 and 1. */
	phase: number
}

export interface Dash extends Point {
	alpha: number
	/** The dash's direction, a unit vector pointing away from the center. */
	toLeft: number
	toTop: number
}

/** The distance between neighboring dots, in CSS pixels. */
export const spacing = 26

/** How far from the center the rings still show. */
export const waveReach = 520

const wavelength = 120
const wavesPerSecond = 0.6
const restingAlpha = 0.18
const crestPush = 8
const jitterShare = 0.6
// The share of cells that hold a dot; the rest stay empty.
const filledShare = 0.5

// A fixed pseudo-random value for a cell, so the grid never reads as a grid
// and the layout is the same on every render.
function cellNoise(left: number, top: number) {
	const value = Math.sin(left * 12.9898 + top * 78.233) * 43758.5453

	return value - Math.floor(value)
}

// Half the cells hold a dot, chosen by the cell's own noise, so the field
// reads as scattered rather than gridded.
export function layoutDots(width: number, height: number): Dot[] {
	const dots: Dot[] = []
	for (let top = spacing / 2; top < height; top += spacing) {
		for (let left = spacing / 2; left < width; left += spacing) {
			const noise = cellNoise(left, top)
			const otherNoise = (noise * 7) % 1
			const keep = (noise * 13) % 1
			if (keep > filledShare) continue
			dots.push({
				left: left + (noise - 0.5) * spacing * jitterShare,
				phase: noise,
				top: top + (otherNoise - 0.5) * spacing * jitterShare,
			})
		}
	}

	return dots
}

// One dot at one moment: rings travel outward from the center, a dot on a
// crest brightens and leans outward, and every dot keeps a faint rest.
export function dashOf(dot: Dot, center: Point, seconds: number): Dash {
	const offsetLeft = dot.left - center.left
	const offsetTop = dot.top - center.top
	const distance = Math.hypot(offsetLeft, offsetTop)
	const toLeft = distance === 0 ? 0 : offsetLeft / distance
	const toTop = distance === 0 ? 0 : offsetTop / distance
	const travel = distance / wavelength - seconds * wavesPerSecond
	const wave = Math.sin((travel + dot.phase * 0.15) * Math.PI * 2)
	const crest = Math.max(0, wave) ** 2
	const fade = Math.max(0, 1 - distance / waveReach)
	const strength = crest * fade

	return {
		alpha: restingAlpha + (1 - restingAlpha) * strength,
		left: dot.left + toLeft * crestPush * strength,
		toLeft,
		toTop,
		top: dot.top + toTop * crestPush * strength,
	}
}
