import { hasPlayed } from './has-played.compute'
import { planPhaseIds, planPhases } from './microruntimes-content.data'
import type { PlanPhaseId } from './microruntimes-content.types'

// The act the session is in: the last one whose opening line has played.
export function phaseOf(playedCount: number): PlanPhaseId {
	const reached = planPhaseIds.filter((phase) =>
		hasPlayed(playedCount, planPhases[phase]),
	)

	return reached.at(-1) ?? 'ticket'
}
