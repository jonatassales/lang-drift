import { hasPlayed } from './has-played.compute'
import { planReveals } from './microruntimes-content.data'
import type { PlanRevealId } from './microruntimes-content.types'

// True once the terminal has completed the line that reveals this part of
// the cards.
export function revealed(playedCount: number, id: PlanRevealId): boolean {
	return hasPlayed(playedCount, planReveals[id])
}
