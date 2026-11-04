'use client'

import { useReducedMotion } from 'motion/react'
import { type ReactNode, useEffect, useRef, useState } from 'react'

import { cn } from 'web/utilities'

import { createTerminalMachine, type Machine } from './create-terminal-machine'

// Adapted from SmoothUI's Mac Terminal (smoothui.dev): a macOS terminal
// window that types its commands and prints their output line by line. It
// plays while it is on screen and the tab is visible, shows the whole session
// at once under reduced motion, and keeps the full transcript available to
// screen readers. The colors are the tokens' roles, so the window follows the
// color scheme, and the traffic lights take the status colors. The window
// takes its height from its className; the transcript fills what the title
// bar leaves and scrolls, so the window never grows as it types.

export type TerminalLineType =
	| 'command'
	| 'output'
	| 'error'
	| 'success'
	| 'comment'
	| 'warning'

export interface TerminalLine {
	/** Milliseconds to wait before the line starts. */
	delay?: number
	id: string
	text: string
	type: TerminalLineType
}

export interface MacTerminalProps {
	/** A control shown at the right of the title bar. */
	action?: ReactNode
	/** The window's classes, among them its height. */
	className?: string
	/** The accessible name of the window. */
	label: string
	lines: TerminalLine[]
	loop?: boolean
	/** Called with the number of completed lines whenever it changes. */
	onProgress?: (completedCount: number) => void
	/** Freezes the session while true. */
	paused?: boolean
	prompt?: string
	title: string
	/** Characters per second. */
	typingSpeed?: number
}

interface TypingDraft {
	chars: number
	line: TerminalLine
}

const TICK_MS = 40
const LOOP_PAUSE_MS = 1500
const OUTPUT_DELAY_MS = 200
const VISIBLE_THRESHOLD = 0.35
const trafficLight = 'size-3 rounded-full'

const lineTones: Record<TerminalLineType, string> = {
	command: 'text-foreground',
	comment: 'text-muted-foreground',
	error: 'text-status-danger',
	output: 'text-foreground/80',
	success: 'text-status-success',
	warning: 'text-status-warning',
}

export function MacTerminal(props: MacTerminalProps) {
	const {
		action,
		className,
		label,
		lines,
		loop = false,
		onProgress,
		paused = false,
		prompt = '~ %',
		title,
		typingSpeed = 24,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const containerRef = useRef<HTMLElement>(null)
	const scrollRef = useRef<HTMLDivElement>(null)
	const machineRef = useRef<Machine>(createTerminalMachine())
	const shouldPlayRef = useRef(false)
	const linesRef = useRef(lines)
	const onProgressRef = useRef(onProgress)
	const reportedCountRef = useRef(0)
	const loopRef = useRef(loop)
	const msPerCharRef = useRef(1000 / typingSpeed)

	linesRef.current = lines
	onProgressRef.current = onProgress
	loopRef.current = loop
	msPerCharRef.current = 1000 / typingSpeed

	const [completedLines, setCompletedLines] = useState<TerminalLine[]>([])
	const [typingDraft, setTypingDraft] = useState<TypingDraft | null>(null)
	const [isInView, setIsInView] = useState(false)
	const [isPageVisible, setIsPageVisible] = useState(true)

	useEffect(() => {
		const node = containerRef.current
		if (!node) return

		const observer = new IntersectionObserver(
			([entry]) => setIsInView(entry?.isIntersecting ?? false),
			{ threshold: VISIBLE_THRESHOLD },
		)
		observer.observe(node)

		return () => observer.disconnect()
	}, [])

	useEffect(() => {
		function updateVisibility() {
			setIsPageVisible(document.visibilityState === 'visible')
		}

		updateVisibility()
		document.addEventListener('visibilitychange', updateVisibility)

		return () =>
			document.removeEventListener('visibilitychange', updateVisibility)
	}, [])

	useEffect(() => {
		shouldPlayRef.current = isInView && isPageVisible && !paused
	}, [isInView, isPageVisible, paused])

	useEffect(() => {
		if (completedLines.length === reportedCountRef.current) return

		reportedCountRef.current = completedLines.length
		onProgressRef.current?.(completedLines.length)
	}, [completedLines])

	useEffect(() => {
		if (!shouldReduceMotion) return

		setCompletedLines(lines)
		setTypingDraft(null)
	}, [shouldReduceMotion, lines])

	useEffect(() => {
		if (shouldReduceMotion) return

		let loopTimeout = 0

		function startNextLine(machine: Machine) {
			machine.lineIndex += 1
			machine.phase = 'delay'
			machine.phaseElapsed = 0
		}

		function tick() {
			const machine = machineRef.current
			if (!shouldPlayRef.current || machine.finished) return

			const currentLines = linesRef.current
			if (machine.lineIndex >= currentLines.length) {
				machine.finished = true
				if (!loopRef.current) return

				loopTimeout = window.setTimeout(() => {
					machineRef.current = createTerminalMachine()
					setCompletedLines([])
					setTypingDraft(null)
				}, LOOP_PAUSE_MS)

				return
			}

			const line = currentLines[machine.lineIndex]
			const isCommand = line.type === 'command'
			const delay = line.delay ?? (isCommand ? 0 : OUTPUT_DELAY_MS)

			machine.phaseElapsed += TICK_MS

			if (machine.phase === 'delay') {
				if (machine.phaseElapsed < delay) return

				if (!isCommand) {
					setCompletedLines((previous) => [...previous, line])
					startNextLine(machine)

					return
				}

				machine.phase = 'typing'
				machine.phaseElapsed = 0
				setTypingDraft({
					chars: 0,
					line,
				})

				return
			}

			const chars = Math.min(
				line.text.length,
				Math.floor(machine.phaseElapsed / msPerCharRef.current),
			)
			setTypingDraft({
				chars,
				line,
			})
			if (chars < line.text.length) return

			setCompletedLines((previous) => [...previous, line])
			setTypingDraft(null)
			startNextLine(machine)
		}

		const interval = window.setInterval(tick, TICK_MS)

		return () => {
			window.clearInterval(interval)
			window.clearTimeout(loopTimeout)
		}
	}, [shouldReduceMotion])

	// biome-ignore lint/correctness/useExhaustiveDependencies: each new line or typed character scrolls the window to its end
	useEffect(() => {
		const scroller = scrollRef.current
		if (!scroller) return

		scroller.scrollTop = scroller.scrollHeight
	}, [completedLines, typingDraft])

	return (
		<figure
			aria-label={label}
			className={cn(
				'm-0 flex flex-col overflow-hidden rounded-xl border border-border bg-card text-left shadow-lg',
				className,
			)}
			ref={containerRef}
		>
			<div className="flex items-center gap-3 border-b border-border bg-muted px-4 py-2.5">
				<div aria-hidden="true" className="flex items-center gap-2">
					<span className={cn(trafficLight, 'bg-status-danger')} />
					<span className={cn(trafficLight, 'bg-status-warning')} />
					<span className={cn(trafficLight, 'bg-status-success')} />
				</div>
				<p className="m-0 flex-1 truncate text-center text-xs font-medium text-muted-foreground">
					{title}
				</p>
				<span className="flex min-w-14 items-center justify-end">{action}</span>
			</div>

			<div
				className="min-h-48 flex-1 contain-size overflow-y-auto px-4 py-3 font-mono text-xs leading-relaxed sm:text-sm"
				ref={scrollRef}
			>
				<pre
					aria-hidden={shouldReduceMotion ? undefined : 'true'}
					className="m-0 font-mono break-words whitespace-pre-wrap"
				>
					{completedLines.map((line) => (
						<div className="flex flex-wrap gap-2" key={line.id}>
							{line.type === 'command' && (
								<span className="text-primary">{prompt}</span>
							)}
							<span className={lineTones[line.type]}>{line.text}</span>
						</div>
					))}
					{typingDraft && (
						<div className="flex flex-wrap gap-2">
							<span className="text-primary">{prompt}</span>
							<span className={lineTones.command}>
								{typingDraft.line.text.slice(0, typingDraft.chars)}
								<span
									aria-hidden="true"
									className="ml-0.5 inline-block h-[1em] w-[0.55em] translate-y-[0.15em] animate-pulse bg-primary align-middle"
								/>
							</span>
						</div>
					)}
				</pre>
				{!shouldReduceMotion && (
					<pre className="sr-only">
						{lines.map((line) => line.text).join('\n')}
					</pre>
				)}
			</div>
			<pre
				aria-hidden="true"
				className="invisible m-0 h-0 overflow-hidden px-4 font-mono text-xs whitespace-pre sm:text-sm"
			>
				{lines.map((line) =>
					line.type === 'command'
						? `${prompt} ${line.text}\n`
						: `${line.text}\n`,
				)}
			</pre>
		</figure>
	)
}
