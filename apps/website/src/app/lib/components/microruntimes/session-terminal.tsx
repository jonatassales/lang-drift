'use client'

import { MacTerminal } from '@langdrift/web/vendors/smoothui/mac-terminal'
import { useAtomValue, useSetAtom } from 'jotai'

import {
	type MicroruntimesMessages,
	planSession,
} from '@data/microruntimes-content'
import { microruntimeSessionAtom, recordPlayedCountAtom } from '@state'

import { SessionControl } from './session-control'

interface SessionTerminalProps {
	control: MicroruntimesMessages['session']['card']['control']
	copy: MicroruntimesMessages['session']['terminal']
}

// The hero's simulated planning session: the terminal plays the script on
// its own clock and records how many lines it has completed in the session,
// which the cards beside it read. Its one control, pause and replay, sits
// in the title bar. A replay remounts it through the run counter. The
// window keeps a fixed height where the cards stack under it, so nothing
// below shifts as it types, and stretches to the cards' height beside
// them, so the two columns end together whatever card is open.
export function SessionTerminal(props: SessionTerminalProps) {
	const { control, copy } = props

	const session = useAtomValue(microruntimeSessionAtom)
	const recordPlayedCount = useSetAtom(recordPlayedCountAtom)
	const lines = planSession.map((line) => ({
		...line,
		text: copy.lines[line.id],
	}))

	return (
		<MacTerminal
			key={session.run}
			action={<SessionControl copy={control} />}
			className="h-[clamp(240px,40vh,360px)] lg:h-auto"
			label={copy.label}
			lines={lines}
			onProgress={recordPlayedCount}
			paused={session.paused}
			prompt="~/project %"
			title={copy.title}
		/>
	)
}
