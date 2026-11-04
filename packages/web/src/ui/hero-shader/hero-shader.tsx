'use client'

import { useEffect, useRef } from 'react'

import { cn } from 'web/utilities'

import type { SurfaceController, SurfaceOptions } from './surface'

export interface HeroShaderProps {
	className?: string
	paused?: boolean
	reducedMotion?: boolean
	/**
	 * `muted` quiets the shader behind content: a faint purple and peach in the
	 * light scheme, grays in the dark one.
	 */
	tone?: 'brand' | 'muted'
}

// A token-only gradient stands behind the canvas: it shows until the program
// draws its first frame, and stays when the browser has no WebGL2.
const fallbackGradient =
	'bg-[radial-gradient(ellipse_at_72%_28%,color-mix(in_srgb,var(--primary)_24%,var(--card)),transparent_60%),linear-gradient(120deg,color-mix(in_srgb,var(--primary-hover)_18%,var(--background)),var(--background)_55%,color-mix(in_srgb,var(--primary)_22%,var(--card)))]'

// The shader reads its colors from the brand roles. The muted tone keeps the
// brand purple in the light scheme and lends its second sheet the rainbow's
// orange, which the pale field turns into peach, all at about a third of the
// opacity, so the waves read as a faint purple-to-peach wash. In the dark
// scheme it swaps the roles for the surface and border grays at under half
// opacity, so the waves barely lift from the black.
const toneClass: Record<NonNullable<HeroShaderProps['tone']>, string> = {
	brand: '',
	muted:
		'opacity-35 [--primary-hover:var(--lg-rainbow-2)] dark:opacity-45 dark:[--primary-hover:var(--border)] dark:[--primary:var(--muted)]',
}

export function HeroShader(props: HeroShaderProps) {
	const {
		className,
		paused = false,
		reducedMotion = false,
		tone = 'brand',
	} = props

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
				toneClass[tone],
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
