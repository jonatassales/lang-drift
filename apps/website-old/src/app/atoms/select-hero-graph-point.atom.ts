import { atom } from 'jotai'

import { heroSelectionAtom } from './hero-selection.atom'

export const selectHeroGraphPointAtom = atom(
	null,
	(get, set, pointIndex: number) => {
		if (get(heroSelectionAtom).pointIndex === pointIndex) return
		set(heroSelectionAtom, {
			pointIndex,
			source: 'graph',
		})
	},
)
