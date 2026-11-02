import type { Rgb } from './palette-types'

export function mix(from: Rgb, to: Rgb, ratio: number): Rgb {
	return [
		from[0] + (to[0] - from[0]) * ratio,
		from[1] + (to[1] - from[1]) * ratio,
		from[2] + (to[2] - from[2]) * ratio,
	]
}
