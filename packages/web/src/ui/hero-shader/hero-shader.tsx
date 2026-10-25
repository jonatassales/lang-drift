'use client'

import { cn } from '@langdrift/web/utilities'
import { useEffect, useRef } from 'react'

import type { SurfaceController, SurfaceOptions } from './surface'

export interface HeroShaderProps {
	className?: string
	paused?: boolean
	reducedMotion?: boolean
}

// A token-only gradient stands behind the canvas: it shows until the program
// draws its first frame, and stays when the browser has no WebGL2.
const fallbackGradient =
	'bg-[radial-gradient(ellipse_at_72%_28%,color-mix(in_srgb,var(--primary)_24%,var(--card)),transparent_60%),linear-gradient(120deg,color-mix(in_srgb,var(--primary-hover)_18%,var(--background)),var(--background)_55%,color-mix(in_srgb,var(--primary)_22%,var(--card)))]'

export function HeroShader(props: HeroShaderProps) {
	const { className, paused = false, reducedMotion = false } = props

	const canvasRef = useRef<HTMLCanvasElement>(null)
	const controllerRef = useRef<SurfaceController | null>(null)
	const optionsRef = useRef<SurfaceOptions>({
		paused,
		reducedMotion,
	})

	useEffect(() => {
		let cancelled = false
		// The program loads after first paint. A decorative failure keeps the
		// gradient, so nothing is reported.
		void import('./surface')
			.then(({ mountSurface }) => {
				const canvas = canvasRef.current
				if (cancelled || !canvas) return
				controllerRef.current = mountSurface(canvas, optionsRef.current) ?? null
			})
			.catch(() => {})

		return () => {
			cancelled = true
			controllerRef.current?.dispose()
			controllerRef.current = null
		}
	}, [])

	useEffect(() => {
		optionsRef.current = {
			paused,
			reducedMotion,
		}
		controllerRef.current?.update(optionsRef.current)
	}, [paused, reducedMotion])

	return (
		<div
			aria-hidden="true"
			className={cn(
				'pointer-events-none absolute inset-0 overflow-hidden bg-background',
				fallbackGradient,
				className,
			)}
		>
			<canvas
				className="block h-full w-full data-[unavailable]:hidden"
				ref={canvasRef}
			/>
		</div>
	)
}
