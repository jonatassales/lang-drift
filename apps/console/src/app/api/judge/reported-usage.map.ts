import type { ProviderUsage } from './judge.types'

// The usage the contract reports: whole, non-negative token counts, so a
// usage the judge refused still goes out in the contract's shape.
export function reportedUsage(usage: ProviderUsage): ProviderUsage {
	return {
		inputTokens: Math.max(0, Math.trunc(usage.inputTokens) || 0),
		outputTokens: Math.max(0, Math.trunc(usage.outputTokens) || 0),
	}
}
