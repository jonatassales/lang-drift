import type { GateFigures, ObservedWorkItem } from './finops-content.data'

export function gateFigures(item: ObservedWorkItem): GateFigures {
	return {
		proposed: item.estimate,
		reserved: item.reserved,
		observed: item.observed,
		released: item.reserved - item.observed,
	}
}
