'use client'

import { useAtom } from 'jotai'

import {
	type MicroruntimesMessages,
	planSession,
} from '@data/microruntimes-content'
import { initialMicroruntimeSession, microruntimeSessionAtom } from '@state'

interface SessionControlProps {
	copy: MicroruntimesMessages['session']['card']['control']
}

// The one control over the hero's moving content: a text button that pauses
// the session while it plays and replays it once it has ended. Replay resets
// the session and bumps the run, which remounts the terminal; the card
// returns to its dashes in the same frame. Visitors who asked for reduced
// motion see the session at rest and no control.
export function SessionControl(props: SessionControlProps) {
	const { copy } = props

	const [session, setSession] = useAtom(microruntimeSessionAtom)
	const finished = session.playedCount >= planSession.length

	const replay = () =>
		setSession({
			...initialMicroruntimeSession,
			run: session.run + 1,
		})
	const togglePause = () =>
		setSession({
			...session,
			paused: !session.paused,
		})

	return (
		<button
			aria-pressed={finished ? undefined : session.paused}
			className="smooth-text-link min-h-11 motion-reduce:hidden"
			onClick={finished ? replay : togglePause}
			type="button"
		>
			{finished ? copy.replay : copy.pause}
		</button>
	)
}
