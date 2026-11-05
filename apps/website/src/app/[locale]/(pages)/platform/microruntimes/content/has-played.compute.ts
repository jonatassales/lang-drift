import type { PlanLineId } from './microruntimes-content.types'
import { planLineIndex } from './plan-line-index.compute'

// True once the terminal has completed the line, so the card can show what
// that line reveals.
export function hasPlayed(playedCount: number, id: PlanLineId): boolean {
	return playedCount > planLineIndex(id)
}
