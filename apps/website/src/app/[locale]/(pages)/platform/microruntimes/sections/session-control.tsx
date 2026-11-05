import { useAtomValue, useSetAtom } from 'jotai'

import type { MicroruntimesMessages } from '@routes/platform/microruntimes/content/microruntimes-content.types'
import {
	microruntimeSessionAtom,
	replaySessionAtom,
	sessionEndedAtom,
	toggleSessionPauseAtom,
} from '@state'

const button =
	'smooth-text-link min-h-11 whitespace-nowrap focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary motion-reduce:hidden'

// The one control over the hero's moving content: a text button that pauses
// and resumes the session while the plan or a run plays, and replays the
// plan once the terminal has ended. Visitors who asked for reduced motion
// see the session at rest and no control.
interface SessionControlProps {
	copy: MicroruntimesMessages['session']['card']['control']
}

export function SessionControl(props: SessionControlProps) {
	const { copy } = props

	const { paused } = useAtomValue(microruntimeSessionAtom)
	const ended = useAtomValue(sessionEndedAtom)
	const replay = useSetAtom(replaySessionAtom)
	const togglePause = useSetAtom(toggleSessionPauseAtom)

	function handleReplay() {
		replay()
	}

	function handleTogglePause() {
		togglePause()
	}

	if (ended) {
		return (
			<button className={button} onClick={handleReplay} type="button">
				{copy.replay}
			</button>
		)
	}

	return (
		<button className={button} onClick={handleTogglePause} type="button">
			{paused ? copy.resume : copy.pause}
		</button>
	)
}
