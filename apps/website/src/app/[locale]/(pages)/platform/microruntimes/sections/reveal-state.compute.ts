import type { RevealValueState } from 'web/ui/reveal-value'

export function revealState(
	valueShown: boolean,
	skeletonShown = false,
): RevealValueState {
	if (valueShown) return 'set'
	if (skeletonShown) return 'resolving'

	return 'unset'
}
