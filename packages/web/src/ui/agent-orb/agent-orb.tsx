'use client'

import type {
	OrbzPresetName,
	OrbzReducedMotion,
	OrbzSize,
	OrbzState,
} from '@neongate-ai/orbz'
import type {} from '@neongate-ai/orbz/react-types'
import { useEffect, useState } from 'react'

export interface AgentOrbProps {
	paused?: boolean
	preset?: OrbzPresetName
	reducedMotion?: OrbzReducedMotion
	size?: OrbzSize
	speed?: number
	state?: OrbzState
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
		void import('@neongate-ai/orbz/browser')
			.then(() => {
				if (!cancelled) setReady(true)
			})
			.catch(() => {})

		return () => {
			cancelled = true
		}
	}, [])

	// Match Orbz's size contract without eagerly importing its browser/runtime graph.
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
					: 'inline-flex shrink-0 rounded-full bg-[radial-gradient(circle_at_35%_30%,color-mix(in_srgb,var(--primary)_35%,white),var(--primary)_50%,color-mix(in_srgb,var(--primary)_60%,black))]'
			}
			style={{
				height: dimension,
				width: dimension,
			}}
		>
			<orb-z
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
