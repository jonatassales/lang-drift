import type { Rgb } from './palette-types'

let probe: CanvasRenderingContext2D | null | undefined

const sentinel = '#010203'
const hexPattern = /^#([0-9a-f]{6})$/i
const rgbPattern = /^rgba?\((\d+),\s*(\d+),\s*(\d+)/

export function normalizeColor(value: string): Rgb | null {
	function probeContext() {
		if (probe === undefined) {
			probe = document.createElement('canvas').getContext('2d')
		}

		return probe
	}

	function channel(hex: string, offset: number) {
		return Number.parseInt(hex.slice(offset, offset + 2), 16) / 255
	}

	const context = probeContext()
	if (!context) return null

	context.fillStyle = sentinel
	context.fillStyle = value
	const serialized = String(context.fillStyle)
	const unparsed =
		serialized === sentinel && value.trim().toLowerCase() !== sentinel
	if (unparsed) return null

	const hex = hexPattern.exec(serialized)?.[1]
	if (hex) return [channel(hex, 0), channel(hex, 2), channel(hex, 4)]

	const rgb = rgbPattern.exec(serialized)
	if (!rgb) return null

	const [, red, green, blue] = rgb
	if (!red || !green || !blue) return null

	return [Number(red) / 255, Number(green) / 255, Number(blue) / 255]
}
