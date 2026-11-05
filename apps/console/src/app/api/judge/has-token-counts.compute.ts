import type { ProviderUsage } from './judge.types'

// Whether the provider reported both token counts as whole, non-negative
// numbers: the only usage the judge can price.
export function hasTokenCounts(usage: ProviderUsage): boolean {
	const counts = [usage.inputTokens, usage.outputTokens]

	return counts.every((count) => Number.isInteger(count) && count >= 0)
}
