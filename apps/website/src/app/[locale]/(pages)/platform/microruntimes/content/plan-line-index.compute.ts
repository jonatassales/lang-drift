import { planSession } from './microruntimes-content.data'
import type { PlanLineId } from './microruntimes-content.types'

export function planLineIndex(id: PlanLineId): number {
	return planSession.findIndex((line) => line.id === id)
}
