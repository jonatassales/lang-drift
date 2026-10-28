import { useAtomValue, useSetAtom } from 'jotai'

import type { MicroruntimesMessages } from '@data/microruntimes-content'
import {
	microruntimeSessionAtom,
	replaySessionAtom,
	sessionEndedAtom,
	toggleSessionPauseAtom,
} from '@state'

interface SessionControlProps {
	copy: MicroruntimesMessages['session']['card']['control']
}

const button = 'smooth-text-link min-h-11 motion-reduce:hidden'

// The one control over the hero's moving content: a text button that pauses
// and resumes the session while the plan or a run plays, and replays the
// plan once the terminal has ended. Visitors who asked for reduced motion
// see the session at rest and no control.
export function SessionControl(props: SessionControlProps) {
	const { copy } = props

	const { paused } = useAtomValue(microruntimeSessionAtom)
	const ended = useAtomValue(sessionEndedAtom)
	const replay = useSetAtom(replaySessionAtom)
	const togglePause = useSetAtom(toggleSessionPauseAtom)

	if (ended) {
		return (
			<button className={button} onClick={() => replay()} type="button">
				{copy.replay}
			</button>
		)
	}

	return (
		<button className={button} onClick={() => togglePause()} type="button">
			{paused ? copy.resume : copy.pause}
		</button>
	)
}
