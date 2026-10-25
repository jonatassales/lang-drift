import { atom } from 'jotai'

import { heroSelectionAtom } from './hero-selection.atom'

export const selectedPointAtom = atom(
	(get) => get(heroSelectionAtom).pointIndex,
	(get, set, pointIndex: number) => {
		if (get(heroSelectionAtom).pointIndex === pointIndex) return
		set(heroSelectionAtom, {
			pointIndex,
			source: 'chart',
		})
	},
)
