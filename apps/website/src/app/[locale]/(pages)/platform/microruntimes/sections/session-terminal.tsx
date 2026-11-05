'use client'

import { useAtomValue, useSetAtom } from 'jotai'
import { createTranslator } from 'next-intl'
import { useEffect, useRef } from 'react'

import { planSession } from '@routes/platform/microruntimes/content/microruntimes-content.data'
import type { MicroruntimesMessages } from '@routes/platform/microruntimes/content/microruntimes-content.types'
import {
	microruntimeSessionAtom,
	recordPlayedCountAtom,
	runScriptAtom,
	startRunAtom,
} from '@state'
import { cn } from 'web/utilities'
import { MacTerminal } from 'web/vendors/smoothui/mac-terminal'

import { SessionControl } from './session-control'

// The one second the Workflow card shows its loader before the run starts.
const launchMs = 1000

// The simulated Microruntimes session: the terminal plays the plan, then,
// once the visitor plays the workflow from the cards, the run, on its own
// clock, and records how many lines it has completed in the session, which
// any cards beside it read. Its one control, pause and replay, sits in the
// title bar. Each run remounts it through the run counter, so the finished
// plan stays on screen through the launch and clears as the run starts.
// Playing the workflow scrolls the page up to the window at once, since it
// only plays while visible and the run starts a second later. The window
// keeps a fixed height, so nothing below shifts as it types; a page that
// sets it beside the cards stretches it to their height with `lg:h-full`,
// so the two columns end together whatever card is open.
interface SessionTerminalProps {
	/** The window's height where it sits, e.g. `lg:h-full` beside the cards. */
	className?: string
	control: MicroruntimesMessages['session']['card']['control']
	copy: MicroruntimesMessages['session']['terminal']
}

export function SessionTerminal(props: SessionTerminalProps) {
	const { className, control, copy } = props

	const windowRef = useRef<HTMLDivElement>(null)

	const session = useAtomValue(microruntimeSessionAtom)
	const script = useAtomValue(runScriptAtom)
	const recordPlayedCount = useSetAtom(recordPlayedCountAtom)
	const startRun = useSetAtom(startRunAtom)

	const { paused, run, stage } = session

	// Launch the workflow run after the visual launch delay.
	useEffect(() => {
		if (stage !== 'launching') return

		const launch = window.setTimeout(() => startRun(), launchMs)

		return () => window.clearTimeout(launch)
	}, [stage, startRun])

	// Scroll terminal into view when a workflow run begins.
	useEffect(() => {
		if (stage !== 'launching') return

		const reducedMotion = window.matchMedia(
			'(prefers-reduced-motion: reduce)',
		).matches

		windowRef.current?.scrollIntoView({
			behavior: reducedMotion ? 'auto' : 'smooth',
			block: 'start',
		})
	}, [stage])

	const runText = createTranslator({
		locale: 'en',
		messages: copy.run,
	})
	const lines =
		stage === 'running'
			? script.map((line) => ({
					delay: line.delay,
					id: line.id,
					text: runText(line.message, line.values),
					type: line.type,
				}))
			: planSession.map((line) => ({
					...line,
					text: copy.lines[line.id],
				}))

	return (
		<div className="min-w-0 scroll-mt-24 lg:h-full" ref={windowRef}>
			<MacTerminal
				key={run}
				action={<SessionControl copy={control} />}
				className={cn('h-[clamp(240px,40vh,360px)]', className)}
				label={copy.label}
				lines={lines}
				onProgress={recordPlayedCount}
				paused={paused}
				prompt="~/project %"
				title={copy.title}
			/>
		</div>
	)
}
