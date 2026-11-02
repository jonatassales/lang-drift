import { mix } from './mix'
import type { RoleColors, Scheme, ShaderPalette } from './palette-types'
import { rotateHue } from './rotate-hue'

const black: [number, number, number] = [0, 0, 0]

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
