export type Rgb = [number, number, number]
type Hsl = [number, number, number]
export type Scheme = 'dark' | 'light'

// The token roles the shader starts from, read from the page at runtime.
export interface RoleColors {
	background: Rgb
	surface: Rgb
	foreground: Rgb
	primary: Rgb
	primaryHover: Rgb
}

// What the shader draws with: the roles it uses directly, the shades derived
// from them, and the scheme they were derived for (0 dark, 1 light).
export interface ShaderPalette {
	background: Rgb
	surface: Rgb
	primary: Rgb
	primaryHover: Rgb
	shadeDeep: Rgb
	shadeSoft: Rgb
	glow: Rgb
	violet: Rgb
	scheme: number
}

const black: Rgb = [0, 0, 0]

function clamp01(value: number) {
	return Math.min(1, Math.max(0, value))
}

export function mix(from: Rgb, to: Rgb, ratio: number): Rgb {
	return [
		from[0] + (to[0] - from[0]) * ratio,
		from[1] + (to[1] - from[1]) * ratio,
		from[2] + (to[2] - from[2]) * ratio,
	]
}

function hueSector(color: Rgb, max: number, chroma: number) {
	const [red, green, blue] = color
	if (max === red) return ((green - blue) / chroma + 6) % 6
	if (max === green) return (blue - red) / chroma + 2

	return (red - green) / chroma + 4
}

function rgbToHsl(color: Rgb): Hsl {
	const [red, green, blue] = color
	const max = Math.max(red, green, blue)
	const min = Math.min(red, green, blue)
	const lightness = (max + min) / 2
	const chroma = max - min
	if (chroma === 0) return [0, 0, lightness]
	const saturation = chroma / (1 - Math.abs(2 * lightness - 1))

	return [hueSector(color, max, chroma) * 60, saturation, lightness]
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

// A neighbor on the color wheel, turned by `degrees` and lifted or lowered by
// `lightnessDelta`, so every extra shade keeps the brand's saturation.
export function rotateHue(color: Rgb, degrees: number, lightnessDelta: number) {
	const [hue, saturation, lightness] = rgbToHsl(color)

	return hslToRgb([
		(hue + degrees + 360) % 360,
		saturation,
		clamp01(lightness + lightnessDelta),
	])
}

// Dark adds light over a near-black field, so its shades lean toward black and
// the hover tone; light lays veils over a pale field, so its shades lean toward
// the background and the foreground ink and sit higher on the lightness scale.
export function derivePalette(
	roles: RoleColors,
	scheme: Scheme,
): ShaderPalette {
	const light = scheme === 'light'

	return {
		background: roles.background,
		surface: roles.surface,
		primary: roles.primary,
		primaryHover: roles.primaryHover,
		shadeDeep: light
			? mix(roles.primary, roles.foreground, 0.35)
			: mix(roles.primary, black, 0.45),
		shadeSoft: mix(roles.primary, roles.background, light ? 0.72 : 0.55),
		glow: light
			? rotateHue(roles.primary, -18, 0.28)
			: rotateHue(roles.primaryHover, -18, 0.1),
		violet: rotateHue(roles.primary, 14, light ? 0.18 : 0.04),
		scheme: light ? 1 : 0,
	}
}

let probe: CanvasRenderingContext2D | null | undefined

// A 2D canvas turns any CSS color the browser knows into `#rrggbb`, or into
// `rgba(...)` when it carries alpha. Without a canvas (jsdom) there is no color.
function probeContext() {
	if (probe === undefined) {
		probe = document.createElement('canvas').getContext('2d')
	}

	return probe
}

const sentinel = '#010203'
const hexPattern = /^#([0-9a-f]{6})$/i
const rgbPattern = /^rgba?\((\d+),\s*(\d+),\s*(\d+)/

function channel(hex: string, offset: number) {
	return Number.parseInt(hex.slice(offset, offset + 2), 16) / 255
}

export function normalizeColor(value: string): Rgb | null {
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

function readRole(styles: CSSStyleDeclaration, variable: string) {
	return normalizeColor(styles.getPropertyValue(variable).trim())
}

// The palette for the element's place in the page: its roles, resolved against
// the color scheme in force on `<html>`.
export function readPalette(element: Element): ShaderPalette | null {
	const styles = getComputedStyle(element)
	const background = readRole(styles, '--background')
	const surface = readRole(styles, '--card')
	const foreground = readRole(styles, '--foreground')
	const primary = readRole(styles, '--primary')
	const primaryHover = readRole(styles, '--primary-hover')
	if (!background || !surface || !foreground || !primary || !primaryHover) {
		return null
	}
	const scheme: Scheme = document.documentElement.classList.contains('dark')
		? 'dark'
		: 'light'

	return derivePalette(
		{
			background,
			surface,
			foreground,
			primary,
			primaryHover,
		},
		scheme,
	)
}
