'use client'

import { cn } from '@langdrift/web/utilities'
import { useReducedMotion } from 'motion/react'
import { type ReactNode, useEffect, useRef } from 'react'

export interface ChromaTextProps {
	children: ReactNode
	className?: string
}

// The fixed rainbow stops, orange through blue and violet, fill the letters
// over a background twice their width.
const chroma =
	'bg-[linear-gradient(100deg,var(--lg-rainbow-2),color-mix(in_srgb,var(--lg-rainbow-2)_55%,var(--lg-rainbow-3)),var(--lg-rainbow-5),var(--lg-rainbow-6),var(--lg-rainbow-2))] bg-size-[220%_100%] bg-clip-text text-transparent'

const flowMs = 7000

// Text the rainbow flows through in one slow, endless loop. Under reduced
// motion, or where the browser cannot animate an element, the gradient holds
// still.
export function ChromaText(props: ChromaTextProps) {
	const { children, className } = props

	const shouldReduceMotion = useReducedMotion()
	const textRef = useRef<HTMLSpanElement>(null)

	useEffect(() => {
		const text = textRef.current
		const canFlow = !shouldReduceMotion && typeof text?.animate === 'function'
		if (!text || !canFlow) return

		const flow = text.animate(
			[{ backgroundPosition: '0% 50%' }, { backgroundPosition: '220% 50%' }],
			{
				duration: flowMs,
				easing: 'linear',
				iterations: Number.POSITIVE_INFINITY,
			},
		)

		return () => flow.cancel()
	}, [shouldReduceMotion])

	return (
		<span className={cn(chroma, className)} ref={textRef}>
			{children}
		</span>
	)
}
