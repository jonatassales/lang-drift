import { useAtomValue, useSetAtom } from 'jotai'

import {
	type MicroruntimesMessages,
	revealed,
} from '@data/microruntimes-content'
import {
	microruntimeSessionAtom,
	replaySessionAtom,
	toggleSessionPauseAtom,
} from '@state'

interface SessionControlProps {
	copy: MicroruntimesMessages['session']['card']['control']
}

const button = 'smooth-text-link min-h-11 motion-reduce:hidden'

// The one control over the hero's moving content: a text button that pauses
// and resumes the session while it plays and replays it once it has ended.
// Visitors who asked for reduced motion see the session at rest and no
// control.
export function SessionControl(props: SessionControlProps) {
	const { copy } = props

	const { paused, playedCount } = useAtomValue(microruntimeSessionAtom)
	const replay = useSetAtom(replaySessionAtom)
	const togglePause = useSetAtom(toggleSessionPauseAtom)

	if (revealed(playedCount, 'finished')) {
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
