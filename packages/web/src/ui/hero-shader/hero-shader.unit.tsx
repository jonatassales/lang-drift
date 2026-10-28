import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { HeroShader } from './hero-shader'
import {
	derivePalette,
	mix,
	type Rgb,
	type RoleColors,
	rotateHue,
} from './palette'

const brand: Rgb = [0.427, 0.039, 0.839]

const darkRoles: RoleColors = {
	background: [0.043, 0.043, 0.047],
	surface: [0.082, 0.082, 0.09],
	foreground: [0.961, 0.961, 0.961],
	primary: brand,
	primaryHover: [0.549, 0.231, 0.937],
}

const lightRoles: RoleColors = {
	background: [0.969, 0.965, 0.949],
	surface: [1, 1, 1],
	foreground: [0.129, 0.129, 0.129],
	primary: brand,
	primaryHover: [0.329, 0.031, 0.647],
}

function brightness(color: Rgb) {
	const [red, green, blue] = color

	return (red + green + blue) / 3
}

describe('HeroShader markup', () => {
	it('renders a hidden, pointer-transparent layer that fills its parent and merges classes', () => {
		const html = renderToString(<HeroShader className="rounded-xl" />)
		const document = new DOMParser().parseFromString(html, 'text/html')
		const layer = document.body.firstElementChild
		const canvas = document.querySelector('canvas')

		expect(layer?.getAttribute('aria-hidden')).toBe('true')
		expect(layer?.classList.contains('absolute')).toBe(true)
		expect(layer?.classList.contains('inset-0')).toBe(true)
		expect(layer?.classList.contains('pointer-events-none')).toBe(true)
		expect(layer?.classList.contains('rounded-xl')).toBe(true)
		expect(canvas).not.toBeNull()
		expect(canvas?.hasAttribute('data-unavailable')).toBe(false)
	})

	it('draws in grays when its tone is monochrome, and in the brand otherwise', () => {
		const neutral = '[--primary:var(--muted-foreground)]'
		const brand = new DOMParser().parseFromString(
			renderToString(<HeroShader />),
			'text/html',
		).body.firstElementChild
		const monochrome = new DOMParser().parseFromString(
			renderToString(<HeroShader tone="monochrome" />),
			'text/html',
		).body.firstElementChild

		expect(monochrome?.classList.contains(neutral)).toBe(true)
		expect(monochrome?.classList.contains('opacity-45')).toBe(true)
		expect(brand?.classList.contains(neutral)).toBe(false)
	})
})

describe('shader palette', () => {
	it('mixes channels linearly', () => {
		const [red, green, blue] = mix([0, 0, 0], [1, 0.5, 0], 0.5)

		expect(red).toBeCloseTo(0.5)
		expect(green).toBeCloseTo(0.25)
		expect(blue).toBeCloseTo(0)
	})

	it('turns a hue around the wheel and lifts its lightness', () => {
		const [red, green, blue] = rotateHue([1, 0, 0], 120, 0)
		const lifted = rotateHue([1, 0, 0], 0, 0.25)

		expect(red).toBeCloseTo(0)
		expect(green).toBeCloseTo(1)
		expect(blue).toBeCloseTo(0)
		expect(brightness(lifted)).toBeGreaterThan(brightness([1, 0, 0]))
	})

	it('derives lighter shades for the light scheme and marks the scheme', () => {
		const dark = derivePalette(darkRoles, 'dark')
		const light = derivePalette(lightRoles, 'light')

		expect(dark.scheme).toBe(0)
		expect(light.scheme).toBe(1)
		expect(brightness(light.shadeSoft)).toBeGreaterThan(
			brightness(dark.shadeSoft),
		)
		expect(brightness(light.glow)).toBeGreaterThan(brightness(dark.glow))
		expect(brightness(light.violet)).toBeGreaterThan(brightness(dark.violet))
		expect(brightness(dark.shadeDeep)).toBeLessThan(brightness(brand))
		expect(dark.primary).toEqual(brand)
	})
})
