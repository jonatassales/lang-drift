'use client'

import { AgentOrb, type AgentOrbState } from '@pack/react/ui/agent-orb'
import { useFormatter, useTranslations } from 'next-intl'
import { useEffect, useRef, useState } from 'react'

import './voice-preview.css'

const scriptedStates: AgentOrbState[] = [
	'listening',
	'thinking',
	'speaking',
	'idle',
]

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
		<div className="voice-preview">
			<div className="voice-orb-wrap" aria-hidden="true">
				<AgentOrb size="190px" speed={0.82} state={state} />
			</div>
			<div className="voice-copy">
				<span className="section-kicker section-kicker-dark">
					{t('kicker')}
				</span>
				<h2>{t('title')}</h2>
				<p>{t('description')}</p>
				<div className="voice-question">
					<span>{t('question')}</span>
					<button
						className="max-w-full whitespace-normal! py-2! [overflow-wrap:anywhere]"
						type="button"
						onClick={runPreview}
					>
						{t('ask')}
					</button>
				</div>
				<output className="website-sr-only">{t(state)}</output>
				<div
					aria-live="polite"
					className={`voice-answer ${showAnswer ? 'voice-answer-visible' : ''}`}
					aria-hidden={!showAnswer}
				>
					<div className="voice-answer-stat">
						<strong data-drift-sign="negative">
							{t('points', { count: format.number(-9) })}
						</strong>
						<span>{t('largest')}</span>
					</div>
					<div>
						<strong>{t('answerTitle')}</strong>
						<p>{t('answer')}</p>
					</div>
				</div>
			</div>
		</div>
	)
}
