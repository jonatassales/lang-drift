import { type PlanLineId, planSession } from './microruntimes-content'

export function planLineIndex(id: PlanLineId): number {
	return planSession.findIndex((line) => line.id === id)
}
