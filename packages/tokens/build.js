import StyleDictionary from 'style-dictionary'
import { getReferences, usesReferences } from 'style-dictionary/utils'

// Each product gets one stylesheet built from the shared tokens and its own
// product theme.
const products = ['vision']

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
		const rule = (selector, tokens) =>
			`${selector} {\n${tokens
				.map((token) => `  --${token.name}: ${cssValue(token)};`)
				.join('\n')}\n}\n`

		const rootTokens = dictionary.allTokens.filter(
			(token) => schemeOf(token) !== 'dark',
		)
		const darkTokens = dictionary.allTokens.filter(
			(token) => schemeOf(token) === 'dark',
		)
		return [
			rule(schemes.get('light'), rootTokens),
			rule(schemes.get('dark'), darkTokens),
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
