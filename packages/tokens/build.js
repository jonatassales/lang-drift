import StyleDictionary from 'style-dictionary'
import { getReferences, usesReferences } from 'style-dictionary/utils'

// Each product gets one stylesheet built from the shared tokens and its own
// product theme.
const products = ['vision', 'finops']

// A token under a color scheme group changes with the scheme. The group names
// the selector it is emitted under and is left out of the variable name.
// Tokens outside any scheme group are emitted under `:root`.
const schemes = new Map([
	['light', ':root'],
	['dark', '.dark'],
])

const schemeOf = (token) => (schemes.has(token.path[0]) ? token.path[0] : null)

StyleDictionary.registerTransform({
	name: 'name/langdrift',
	type: 'name',
	transform: (token) =>
		(schemeOf(token) ? token.path.slice(1) : token.path).join('-'),
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
		const light = tokens.filter((token) => schemeOf(token) === 'light')
		const dark = tokens.filter((token) => schemeOf(token) === 'dark')
		// Roles and status colors, named without their prefix and group.
		const namesIn = (group) =>
			light
				.filter((token) => token.path.slice(1, 3).join('.') === group)
				.map((token) => token.path.slice(3).join('-'))
		const roles = namesIn('lg.color')
		const statuses = namesIn('lg.status')
		// Other tokens name the Tailwind theme variable they set.
		const themed = tokens.filter(
			(token) => token.$extensions?.['dev.langdrift.tailwind'],
		)

		return [
			rule(schemes.get('light'), [
				...tokens.filter((token) => !schemeOf(token)).map(variable),
				...light.map(variable),
			]),
			rule(schemes.get('dark'), dark.map(variable)),
			// Declared under both selectors so the references resolve against the
			// scheme of the element, including a `.dark` subtree.
			`/* Shadcn and SmoothUI components read these variables. */\n${rule(
				':root, .dark',
				roles.map((role) => [`--${role}`, `var(--lg-color-${role})`]),
			)}`,
			`/* Tailwind keeps its layout scales; only the tokens set its colors. */\n${rule(
				'@theme inline',
				[
					['--color-*', 'initial'],
					...roles.map((role) => [`--color-${role}`, `var(--${role})`]),
					...statuses.map((status) => [
						`--color-status-${status}`,
						`var(--lg-status-${status})`,
					]),
					...themed.map((token) => [
						token.$extensions['dev.langdrift.tailwind'],
						`var(--${token.name})`,
					]),
				],
			)}`,
		].join('\n')
	},
})

for (const product of products) {
	const tokens = new StyleDictionary({
		source: ['shared/*.json', `themes/${product}.json`],
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
						destination: `${product}.css`,
						format: 'css/langdrift',
					},
				],
			},
		},
	})
	await tokens.buildAllPlatforms()
}
