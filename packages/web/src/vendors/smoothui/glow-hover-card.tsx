'use client'

import { cn } from '@langdrift/web/utilities'
import { useReducedMotion } from 'motion/react'
import type { CSSProperties, ReactNode } from 'react'
import { useEffect, useRef } from 'react'

// Adapted from SmoothUI's Glow Hover Card (smoothui.dev): a grid of cards
// with a glowing twin of the same grid laid over it, revealed only inside a
// soft circle that follows the pointer. The twin shares the grid's classes,
// so its cards sit exactly over the originals without measuring them; it is
// inert, so its copies of any link are neither focusable nor announced. The
// glow takes the brand purple from `--primary`.

export interface GlowHoverCardItem {
	content: ReactNode
	id: string
}

export interface GlowHoverCardsProps {
	cardClassName?: string
	className?: string
	items: GlowHoverCardItem[]
	maskSize?: number
}

const mask =
	'radial-gradient(var(--glow-size) var(--glow-size) at var(--glow-x) var(--glow-y), #000 1%, transparent 50%)'

export function GlowHoverCards(props: GlowHoverCardsProps) {
	const { cardClassName, className, items, maskSize = 400 } = props

	const containerRef = useRef<HTMLDivElement>(null)
	const overlayRef = useRef<HTMLUListElement>(null)
	const shouldReduceMotion = useReducedMotion()

	// The pointer is written straight to the overlay's custom properties, so
	// following it never re-renders the cards.
	useEffect(() => {
		const container = containerRef.current
		const overlay = overlayRef.current
		if (!container || !overlay || shouldReduceMotion) return

		const follow = (event: PointerEvent) => {
			const bounds = container.getBoundingClientRect()
			overlay.style.setProperty('--glow-x', `${event.clientX - bounds.left}px`)
			overlay.style.setProperty('--glow-y', `${event.clientY - bounds.top}px`)
			overlay.style.opacity = '1'
		}
		const fade = () => {
			overlay.style.opacity = '0'
		}

		container.addEventListener('pointermove', follow)
		container.addEventListener('pointerleave', fade)

		return () => {
			container.removeEventListener('pointermove', follow)
			container.removeEventListener('pointerleave', fade)
		}
	}, [shouldReduceMotion])

	return (
		<div className="relative" ref={containerRef}>
			<ul className={cn('m-0 list-none p-0', className)}>
				{items.map((item) => (
					<li className={cardClassName} key={item.id}>
						{item.content}
					</li>
				))}
			</ul>
			{!shouldReduceMotion && (
				<ul
					aria-hidden="true"
					className={cn(
						'pointer-events-none absolute inset-0 m-0 list-none p-0 opacity-0 transition-opacity duration-200 select-none',
						className,
					)}
					inert
					ref={overlayRef}
					style={
						{
							'--glow-size': `${maskSize}px`,
							maskImage: mask,
							WebkitMaskImage: mask,
						} as CSSProperties
					}
				>
					{items.map((item) => (
						<li
							className={cn(
								cardClassName,
								'border-primary bg-primary/15 shadow-lg ring-1 shadow-primary/15 ring-primary ring-inset',
							)}
							key={item.id}
						>
							{item.content}
						</li>
					))}
				</ul>
			)}
		</div>
	)
}
