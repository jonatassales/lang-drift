import { derivePalette } from './derive-palette'
import { normalizeColor } from './normalize-color'
import type { Scheme, ShaderPalette } from './palette-types'

// The palette for the element's place in the page: its roles, resolved against
// the color scheme in force on `<html>`.
export function readPalette(element: Element): ShaderPalette | null {
	function readRole(styles: CSSStyleDeclaration, variable: string) {
		return normalizeColor(styles.getPropertyValue(variable).trim())
	}

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
