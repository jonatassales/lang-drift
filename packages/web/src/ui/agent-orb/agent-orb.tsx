'use client'

import type {
	OrboPresetName,
	OrboReducedMotion,
	OrboSize,
	OrboState,
} from 'orbo-voice'
import type {} from 'orbo-voice/react-types'
import { useEffect, useState } from 'react'

export interface AgentOrbProps {
	paused?: boolean
	preset?: OrboPresetName
	reducedMotion?: OrboReducedMotion
	size?: OrboSize
	speed?: number
	state?: OrboState
}

export function AgentOrb(props: AgentOrbProps) {
	const {
		paused = false,
		preset = 'peach',
		reducedMotion = 'system',
		size = '160px',
		speed = 0.9,
		state = 'idle',
	} = props

	const [ready, setReady] = useState(false)

	useEffect(() => {
		let cancelled = false
		// Decorative runtime failure retains the size-preserving fallback.
		void import('orbo-voice/browser')
			.then(() => {
				if (!cancelled) setReady(true)
			})
			.catch(() => {})

		return () => {
			cancelled = true
		}
	}, [])

	// Match Orbo's size contract without eagerly importing its browser/runtime graph.
	const rawSize = String(size).trim()
	const numericSize = Number(rawSize)
	const dimension =
		!rawSize ||
		(typeof size === 'number' && !Number.isFinite(size)) ||
		(Number.isFinite(numericSize) && numericSize <= 0)
			? '16rem'
			: Number.isFinite(numericSize)
				? `${numericSize}px`
				: rawSize

	return (
		<span
			aria-hidden="true"
			className={
				ready
					? 'inline-flex shrink-0'
					: 'inline-flex shrink-0 rounded-full bg-radial-[circle_at_35%_30%] from-primary-highlight via-primary via-50% to-primary-shade'
			}
			style={{
				height: dimension,
				width: dimension,
			}}
		>
			<orb-o
				paused={paused}
				preset={preset}
				reduced-motion={reducedMotion}
				size={dimension}
				speed={speed}
				state={state}
			/>
		</span>
	)
}
