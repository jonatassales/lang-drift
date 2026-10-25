'use client'

import { MacTerminal } from '@langdrift/web/vendors/smoothui/mac-terminal'
import { useAtom } from 'jotai'

import {
	type MicroruntimesMessages,
	planSession,
} from '@data/microruntimes-content'
import { microruntimeSessionAtom } from '@state'

interface SessionTerminalProps {
	copy: MicroruntimesMessages['session']['terminal']
}

// The hero's simulated planning session: the terminal plays the script on
// its own clock and writes how many lines it has completed to the session,
// which the ticket card beside it reads. A replay remounts it through the
// run counter. Its height follows the viewport so the ticket head fits in a
// phone's first screen.
export function SessionTerminal(props: SessionTerminalProps) {
	const { copy } = props

	const [session, setSession] = useAtom(microruntimeSessionAtom)
	const lines = planSession.map((line) => ({
		...line,
		text: copy.lines[line.id],
	}))

	return (
		<MacTerminal
			key={session.run}
			label={copy.label}
			lines={lines}
			maxHeight="clamp(200px, 38vh, 320px)"
			onProgress={(playedCount) =>
				setSession((current) => ({
					...current,
					playedCount,
				}))
			}
			paused={session.paused}
			prompt="~/shop %"
			title={copy.title}
		/>
	)
}
