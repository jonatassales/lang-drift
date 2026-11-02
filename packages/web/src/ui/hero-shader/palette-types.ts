export type Rgb = [number, number, number]
export type Hsl = [number, number, number]
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
