import { reserveMarginPercent } from './microruntimes-workflow.data'

export function reserveCents(estimateCents: number): number {
	return Math.ceil((estimateCents * (100 + reserveMarginPercent)) / 100)
}
