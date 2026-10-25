'use client'

import { motion, useReducedMotion } from 'motion/react'
import type { ReactNode } from 'react'

// Adapted from SmoothUI's Short Slide Down (smoothui.dev): each line drops in
// from above, slightly blurred, one after another, until the stack locks in
// place. Here the lines are a list that builds once and stays, so every
// item keeps its place in the layout from the first paint and no link ever
// leaves the page.

export interface ShortSlideDownItem {
	content: ReactNode
	id: string
}

export interface ShortSlideDownProps {
	className?: string
	/** Seconds between one line and the next. */
	interval?: number
	items: ShortSlideDownItem[]
	label?: string
}

const BUILD_EASE = [0.2, 0.8, 0.2, 1] as const

export function ShortSlideDown(props: ShortSlideDownProps) {
	const { className, interval = 0.5, items, label } = props

	const shouldReduceMotion = useReducedMotion()

	return (
		<ul aria-label={label} className={className}>
			{items.map((item, index) => (
				<motion.li
					animate={{
						filter: 'blur(0px)',
						opacity: 1,
						scale: 1,
						y: 0,
					}}
					initial={
						shouldReduceMotion
							? false
							: {
									filter: 'blur(2.4px)',
									opacity: 0,
									scale: 0.992,
									y: -28,
								}
					}
					key={item.id}
					transition={{
						delay: index * interval,
						duration: index === 0 ? 0.36 : 0.5,
						ease: BUILD_EASE,
					}}
				>
					{item.content}
				</motion.li>
			))}
		</ul>
	)
}
