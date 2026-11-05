import { runHasPlayed } from '@routes/platform/microruntimes/content/run-has-played.compute'

import type { RunOutcome } from './run-outcome.compute'

export function runEnd(outcome: RunOutcome) {
	const { playedCount, script } = outcome
	const played = (message: Parameters<typeof runHasPlayed>[2]) =>
		runHasPlayed(script, playedCount, message)

	if (played('gatePaused')) return 'paused'
	if (played('readyForReview')) return 'readyForReview'
	if (played('readyToMerge')) return 'readyToMerge'

	return null
}
