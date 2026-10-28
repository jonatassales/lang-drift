'use client'

import { AgentOrb, type AgentOrbState } from '@langdrift/web/ui/agent-orb'
import { cn } from '@langdrift/web/utilities'
import { useFormatter, useTranslations } from 'next-intl'
import { useEffect, useRef, useState } from 'react'

const scriptedStates: AgentOrbState[] = [
	'listening',
	'thinking',
	'speaking',
	'idle',
]

// The orb beside the sample question and its scripted answer. The section
// around it says what the preview is; this component only plays it.
export function VoicePreview() {
	const t = useTranslations('home.voice')
	const format = useFormatter()
	const [state, setState] = useState<AgentOrbState>('idle')
	const [showAnswer, setShowAnswer] = useState(false)
	const timers = useRef<Array<ReturnType<typeof setTimeout>>>([])

	useEffect(() => {
		return () => {
			for (const timer of timers.current) clearTimeout(timer)
		}
	}, [])

	function runPreview() {
		for (const timer of timers.current) clearTimeout(timer)
		timers.current = []
		setShowAnswer(false)

		scriptedStates.forEach((nextState, index) => {
			const timer = setTimeout(() => {
				setState(nextState)
				if (nextState === 'speaking') setShowAnswer(true)
			}, index * 850)
			timers.current.push(timer)
		})
	}

	return (
		<div className="grid w-full items-center gap-7 md:grid-cols-[minmax(220px,0.75fr)_minmax(0,1.25fr)] md:gap-16">
			<div
				aria-hidden="true"
				className="flex min-h-[190px] items-center justify-center md:min-h-[250px]"
			>
				<AgentOrb size="190px" speed={0.82} state={state} />
			</div>
			<div className="grid gap-[18px]">
				<div className="flex flex-col items-stretch justify-between gap-[18px] border-t border-border pt-[18px] md:flex-row md:items-center">
					<span className="text-[13px]">{t('question')}</span>
					<button
						className="min-h-[38px] max-w-full shrink-0 cursor-pointer rounded-lg border-0 bg-primary px-3.5 py-2 text-primary-foreground whitespace-normal [overflow-wrap:anywhere]"
						type="button"
						onClick={runPreview}
					>
						{t('ask')}
					</button>
				</div>
				<output className="sr-only">{t(state)}</output>
				<div
					aria-live="polite"
					className={cn(
						'grid items-start gap-[18px] rounded-xl border border-border p-4 transition-[opacity,transform] duration-200 ease-out motion-reduce:transition-none md:grid-cols-[auto_1fr]',
						showAnswer
							? 'translate-y-0 opacity-100'
							: 'translate-y-2 opacity-0',
					)}
					aria-hidden={!showAnswer}
				>
					<div className="grid gap-[3px]">
						<strong className="text-[22px]" data-drift-sign="negative">
							{t('points', { count: format.number(-9) })}
						</strong>
						<span className="text-[10px] text-muted-foreground">
							{t('largest')}
						</span>
					</div>
					<div>
						<strong>{t('answerTitle')}</strong>
						<p className="mt-[5px] mb-0 text-[10px] text-muted-foreground">
							{t('answer')}
						</p>
					</div>
				</div>
			</div>
		</div>
	)
}
