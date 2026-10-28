import { atom } from 'jotai'

import { type PlanPhaseId, phaseOf } from '@data/microruntimes-content'

// The playback of the Microruntimes hero's simulated planning session: how
// many terminal lines have completed, whether the visitor paused it, a run
// counter that remounts the terminal on replay, and the card the visitor
// opened by hand, if any. The terminal writes the count from its own clock;
// the cards derive every field from it.
export interface MicroruntimeSession {
	/** The card the visitor opened; null follows the session's act. */
	openedCard: PlanPhaseId | null
	playedCount: number
	paused: boolean
	run: number
}

export const initialMicroruntimeSession: MicroruntimeSession = {
	openedCard: null,
	playedCount: 0,
	paused: false,
	run: 0,
}

export const microruntimeSessionAtom = atom<MicroruntimeSession>(
	initialMicroruntimeSession,
)

// The card open beside the terminal: the one the visitor chose, else the
// one the session's act is filling.
export const openCardAtom = atom((get) => {
	const { openedCard, playedCount } = get(microruntimeSessionAtom)

	return openedCard ?? phaseOf(playedCount)
})

// Records how many lines the terminal has completed. When that moves the
// session into its next act, the act takes the cards back from the visitor.
export const recordPlayedCountAtom = atom(
	null,
	(get, set, playedCount: number) => {
		const session = get(microruntimeSessionAtom)
		const actChanged = phaseOf(playedCount) !== phaseOf(session.playedCount)

		set(microruntimeSessionAtom, {
			...session,
			openedCard: actChanged ? null : session.openedCard,
			playedCount,
		})
	},
)

// Opens a card by hand. The open card stays open: one card is always
// expanded, so the column beside the terminal never collapses to its
// headers.
export const openCardByHandAtom = atom(null, (get, set, card: PlanPhaseId) => {
	const session = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...session,
		openedCard: card,
	})
})

// Replays the session from the start: the count, the pause and the opened
// card reset, and the run bumps so the terminal remounts with a fresh
// machine.
export const replaySessionAtom = atom(null, (get, set) => {
	const { run } = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...initialMicroruntimeSession,
		run: run + 1,
	})
})

// Pauses the session where it is, or resumes it.
export const toggleSessionPauseAtom = atom(null, (get, set) => {
	const session = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...session,
		paused: !session.paused,
	})
})
