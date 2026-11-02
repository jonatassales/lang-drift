import type { Hsl, Rgb } from './palette-types'

// A neighbor on the color wheel, turned by `degrees` and lifted or lowered by
// `lightnessDelta`, so every extra shade keeps the brand's saturation.
export function rotateHue(
	color: Rgb,
	degrees: number,
	lightnessDelta: number,
): Rgb {
	function clamp01(value: number) {
		return Math.min(1, Math.max(0, value))
	}

	function hueSector(source: Rgb, max: number, chroma: number) {
		const [red, green, blue] = source
		if (max === red) return ((green - blue) / chroma + 6) % 6
		if (max === green) return (blue - red) / chroma + 2

		return (red - green) / chroma + 4
	}

	function rgbToHsl(source: Rgb): Hsl {
		const [red, green, blue] = source
		const max = Math.max(red, green, blue)
		const min = Math.min(red, green, blue)
		const lightness = (max + min) / 2
		const chroma = max - min
		if (chroma === 0) return [0, 0, lightness]
		const saturation = chroma / (1 - Math.abs(2 * lightness - 1))

		return [hueSector(source, max, chroma) * 60, saturation, lightness]
	}

	function sectorColor(hue: number, chroma: number, secondary: number): Rgb {
		switch (Math.floor(hue / 60) % 6) {
			case 0:
				return [chroma, secondary, 0]
			case 1:
				return [secondary, chroma, 0]
			case 2:
				return [0, chroma, secondary]
			case 3:
				return [0, secondary, chroma]
			case 4:
				return [secondary, 0, chroma]
			default:
				return [chroma, 0, secondary]
		}
	}

	function hslToRgb(hsl: Hsl): Rgb {
		const [hue, saturation, lightness] = hsl
		const chroma = (1 - Math.abs(2 * lightness - 1)) * saturation
		const secondary = chroma * (1 - Math.abs(((hue / 60) % 2) - 1))
		const match = lightness - chroma / 2
		const [red, green, blue] = sectorColor(hue, chroma, secondary)

		return [red + match, green + match, blue + match]
	}

	const [hue, saturation, lightness] = rgbToHsl(color)

	return hslToRgb([
		(hue + degrees + 360) % 360,
		saturation,
		clamp01(lightness + lightnessDelta),
	])
}
