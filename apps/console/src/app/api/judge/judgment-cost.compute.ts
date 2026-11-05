import type { JudgeCost, ModelPrices, ProviderUsage } from './judge.types'

// What a usage costs at a model's list prices: input and output tokens each
// at their own rate per million, the sum rounded up to the cent, so a single
// token costs a cent and no used token is free. The figure is the console's
// and illustrative until the owner prices a plan.
export function judgmentCost(
	usage: ProviderUsage,
	prices: ModelPrices,
): JudgeCost {
	const inputMillionths = usage.inputTokens * prices.inputCentsPerMillion
	const outputMillionths = usage.outputTokens * prices.outputCentsPerMillion

	return {
		cents: Math.ceil((inputMillionths + outputMillionths) / 1_000_000),
		illustrative: true,
		provenance: 'console',
	}
}
