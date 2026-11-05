import { type PlanLineId, planSession } from './microruntimes-content.data'

export function planLineIndex(id: PlanLineId): number {
	return planSession.findIndex((line) => line.id === id)
}
