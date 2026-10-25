import { atom } from 'jotai'

export interface HeroSelection {
	pointIndex: number
	source: 'initial' | 'chart' | 'graph'
}

export const heroSelectionAtom = atom<HeroSelection>({
	pointIndex: 2,
	source: 'initial',
})
