import type { ModelPrices } from './judge.types'

// The list prices of the models the console can offer, in cents per million
// tokens, uncached. A model the environment offers must have a price here.
export const modelPrices: Readonly<Partial<Record<string, ModelPrices>>> = {
	'claude-opus-5-5': {
		inputCentsPerMillion: 400,
		outputCentsPerMillion: 2000,
	},
	'claude-sonnet-5-5': {
		inputCentsPerMillion: 200,
		outputCentsPerMillion: 1000,
	},
	'claude-haiku-5-5': {
		inputCentsPerMillion: 10,
		outputCentsPerMillion: 50,
	},
}
