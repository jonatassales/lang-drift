export type CardValueState = 'unset' | 'resolving' | 'set'

export function revealState(
	valueShown: boolean,
	skeletonShown = false,
): CardValueState {
	if (valueShown) return 'set'
	if (skeletonShown) return 'resolving'

	return 'unset'
}
