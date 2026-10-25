'use client'

import { cn } from '@langdrift/web/utilities'
import { useReducedMotion } from 'motion/react'
import { useEffect, useRef, useState } from 'react'

// Adapted from SmoothUI's Mac Terminal (smoothui.dev): a macOS terminal
// window that types its commands and prints their output line by line. It
// plays while it is on screen and the tab is visible, shows the whole session
// at once under reduced motion, and keeps the full transcript available to
// screen readers. The colors are the tokens' roles, so the window follows the
// color scheme, and the traffic lights take the status colors.

export type TerminalLineType =
	| 'command'
	| 'output'
	| 'error'
	| 'success'
	| 'comment'

export interface TerminalLine {
	/** Milliseconds to wait before the line starts. */
	delay?: number
	id: string
	text: string
	type: TerminalLineType
}

export interface MacTerminalProps {
	className?: string
	/** The accessible name of the window. */
	label: string
	lines: TerminalLine[]
	loop?: boolean
	maxHeight?: number | string
	prompt?: string
	title: string
	/** Characters per second. */
	typingSpeed?: number
}

type MachinePhase = 'delay' | 'typing'

interface Machine {
	finished: boolean
	lineIndex: number
	phase: MachinePhase
	phaseElapsed: number
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
}

const createMachine = (): Machine => ({
	finished: false,
	lineIndex: 0,
	phase: 'delay',
	phaseElapsed: 0,
})

export function MacTerminal(props: MacTerminalProps) {
	const {
		className,
		label,
		lines,
		loop = false,
		maxHeight = 360,
		prompt = '~ %',
		title,
		typingSpeed = 24,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const containerRef = useRef<HTMLElement>(null)
	const scrollRef = useRef<HTMLDivElement>(null)
	const machineRef = useRef<Machine>(createMachine())
	const shouldPlayRef = useRef(false)
	const linesRef = useRef(lines)
	const loopRef = useRef(loop)
	const msPerCharRef = useRef(1000 / typingSpeed)

	linesRef.current = lines
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
		const updateVisibility = () =>
			setIsPageVisible(document.visibilityState === 'visible')

		updateVisibility()
		document.addEventListener('visibilitychange', updateVisibility)

		return () =>
			document.removeEventListener('visibilitychange', updateVisibility)
	}, [])

	useEffect(() => {
		shouldPlayRef.current = isInView && isPageVisible
	}, [isInView, isPageVisible])

	useEffect(() => {
		if (!shouldReduceMotion) return

		setCompletedLines(lines)
		setTypingDraft(null)
	}, [shouldReduceMotion, lines])

	useEffect(() => {
		if (shouldReduceMotion) return

		let loopTimeout = 0

		const startNextLine = (machine: Machine) => {
			machine.lineIndex += 1
			machine.phase = 'delay'
			machine.phaseElapsed = 0
		}

		const tick = () => {
			const machine = machineRef.current
			if (!shouldPlayRef.current || machine.finished) return

			const currentLines = linesRef.current
			if (machine.lineIndex >= currentLines.length) {
				machine.finished = true
				if (!loopRef.current) return

				loopTimeout = window.setTimeout(() => {
					machineRef.current = createMachine()
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
				'm-0 overflow-hidden rounded-xl border border-border bg-card text-left shadow-lg',
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
				<span aria-hidden="true" className="w-14" />
			</div>

			<div
				className="overflow-y-auto px-4 py-3 font-mono text-xs leading-relaxed sm:text-sm"
				ref={scrollRef}
				style={{
					maxHeight:
						typeof maxHeight === 'number' ? `${maxHeight}px` : maxHeight,
				}}
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
				<pre
					aria-hidden="true"
					className="invisible m-0 h-0 overflow-hidden font-mono whitespace-pre"
				>
					{lines.map((line) =>
						line.type === 'command'
							? `${prompt} ${line.text}\n`
							: `${line.text}\n`,
					)}
				</pre>
				{!shouldReduceMotion && (
					<pre className="sr-only">
						{lines.map((line) => line.text).join('\n')}
					</pre>
				)}
			</div>
		</figure>
	)
}
