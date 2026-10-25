import { atom } from 'jotai'

// The playback of the Microruntimes hero's simulated planning session: how
// many terminal lines have completed, whether the visitor paused it, and a
// run counter that remounts the terminal on replay. The terminal writes the
// count from its own clock; the ticket card derives every field from it.
export interface MicroruntimeSession {
	playedCount: number
	paused: boolean
	run: number
}

export const initialMicroruntimeSession: MicroruntimeSession = {
	playedCount: 0,
	paused: false,
	run: 0,
}

export const microruntimeSessionAtom = atom<MicroruntimeSession>(
	initialMicroruntimeSession,
)

// Replays the session from the start: the count and the pause reset, and
// the run bumps so the terminal remounts with a fresh machine.
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
