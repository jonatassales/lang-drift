import { hasPlayed } from './has-played.compute'
import { type PlanRevealId, planReveals } from './microruntimes-content'

// True once the terminal has completed the line that reveals this part of
// the cards.
export function revealed(playedCount: number, id: PlanRevealId): boolean {
	return hasPlayed(playedCount, planReveals[id])
}
