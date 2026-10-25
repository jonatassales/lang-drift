import StyleDictionary from 'style-dictionary'
import { getReferences, usesReferences } from 'style-dictionary/utils'

// Each theme gets one stylesheet built from the shared tokens and the theme.
const themes = ['website']

// A token under a color scheme group changes with the scheme. The group names
// the selector it is emitted under and is left out of the variable name.
// Tokens outside any scheme group are emitted under `:root`.
const schemes = new Map([
	['light', ':root'],
	['dark', '.dark'],
])

// A token naming a Tailwind theme variable here sets that variable.
const tailwindExtension = 'dev.langdrift.tailwind'

const schemeOf = (token) => (schemes.has(token.path[0]) ? token.path[0] : null)

// A token's path without its scheme group.
const pathOf = (token) => (schemeOf(token) ? token.path.slice(1) : token.path)

StyleDictionary.registerTransform({
	name: 'name/langdrift',
	type: 'name',
	transform: (token) => pathOf(token).join('-'),
})

StyleDictionary.registerFormat({
	name: 'css/langdrift',
	format: ({ dictionary }) => {
		// An alias stays a reference, so changing one token updates every token
		// that points to it.
		const cssValue = (token) => {
			if (!usesReferences(token.original.$value)) return token.$value
			const [target] = getReferences(token.original.$value, dictionary.tokens, {
				usesDtcg: true,
			})
			return `var(--${target.name})`
		}
		const variable = (token) => [`--${token.name}`, cssValue(token)]
		const rule = (selector, declarations) =>
			`${selector} {\n${declarations
				.map(([name, value]) => `  ${name}: ${value};`)
				.join('\n')}\n}\n`

		const tokens = dictionary.allTokens
		const schemeless = tokens.filter((token) => !schemeOf(token))
		const light = tokens.filter((token) => schemeOf(token) === 'light')
		const dark = tokens.filter((token) => schemeOf(token) === 'dark')
		// The tokens of a group, named without their prefix and group. Dark
		// shares its names with light, so light stands for both schemes.
		const namesInGroup = (group, candidates) =>
			candidates
				.map(pathOf)
				.filter((path) => path.slice(0, 2).join('.') === group)
				.map((path) => path.slice(2).join('-'))
		const roles = namesInGroup('lg.color', light)
		// Colors that never take the brand, exposed to Tailwind by group.
		const fixedGroups = ['status', 'classification', 'drift']
		const tailwindTokens = tokens.filter(
			(token) => token.$extensions?.[tailwindExtension],
		)

		return [
			rule(schemes.get('light'), [
				...schemeless.map(variable),
				...light.map(variable),
			]),
			rule(schemes.get('dark'), dark.map(variable)),
			// Declared under both selectors so the references resolve against the
			// scheme of the element, including a `.dark` subtree.
			`/* Shadcn and SmoothUI components read these variables. */\n${rule(
				[...schemes.values()].join(', '),
				roles.map((role) => [`--${role}`, `var(--lg-color-${role})`]),
			)}`,
			"/* Tailwind's `dark:` variant follows the `.dark` class. */\n@custom-variant dark (&:where(.dark, .dark *));\n",
			`/* Tailwind keeps its layout scales; only the tokens set its colors. */\n${rule(
				'@theme inline',
				[
					['--color-*', 'initial'],
					...roles.map((role) => [`--color-${role}`, `var(--${role})`]),
					...fixedGroups.flatMap((group) =>
						namesInGroup(`lg.${group}`, [...schemeless, ...light]).map(
							(name) => [
								`--color-${group}-${name}`,
								`var(--lg-${group}-${name})`,
							],
						),
					),
					...tailwindTokens.map((token) => [
						token.$extensions[tailwindExtension],
						`var(--${token.name})`,
					]),
				],
			)}`,
		].join('\n')
	},
})

for (const theme of themes) {
	const styleDictionary = new StyleDictionary({
		source: ['shared/*.json', `themes/${theme}.json`],
		usesDtcg: true,
		// Every scheme-dependent token shares its name with its counterpart in the
		// other scheme, which Style Dictionary would report as a collision.
		// Verbose logging makes a broken reference name the tokens involved.
		log: {
			warnings: 'disabled',
			verbosity: 'verbose',
		},
		platforms: {
			css: {
				transforms: ['name/langdrift'],
				buildPath: 'dist/',
				files: [
					{
						destination: `${theme}.css`,
						format: 'css/langdrift',
					},
				],
			},
		},
	})
	await styleDictionary.buildAllPlatforms()
}
