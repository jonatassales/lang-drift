import { act } from '@testing-library/react'
import type { createStore } from 'jotai'

import { recordPlayedCountAtom } from '@state'

// Plays the launched run to a line, the way the terminal records it.
export function playRunTo(
	store: ReturnType<typeof createStore>,
	count: number,
) {
	act(() => {
		store.set(recordPlayedCountAtom, count)
	})
}
