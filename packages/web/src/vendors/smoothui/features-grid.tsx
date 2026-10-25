'use client'

import { cn } from '@langdrift/web/utilities'
import { motion, useReducedMotion } from 'motion/react'
import type { ReactNode } from 'react'
import { useEffect, useState } from 'react'

// Adapted from SmoothUI's Features Grid block (smoothui.dev): a heading and a
// grid of cards that rise into view one after another and lift on hover. The
// content comes from the caller.

export interface FeaturesGridItem {
	description: string
	icon: ReactNode
	id: string
	// A string, or a rendered link: give the link `after:absolute
	// after:inset-0` and the whole card becomes its target.
	title: ReactNode
}

export interface FeaturesGridProps {
	className?: string
	headingId: string
	items: FeaturesGridItem[]
	lead: string
	title: string
}

const SPRING = {
	bounce: 0.1,
	duration: 0.25,
	type: 'spring' as const,
}

const STAGGER = 0.05

export function FeaturesGrid(props: FeaturesGridProps) {
	const { className, headingId, items, lead, title } = props

	const shouldReduceMotion = useReducedMotion()
	const [isHoverDevice, setIsHoverDevice] = useState(false)

	useEffect(() => {
		const query = window.matchMedia('(hover: hover) and (pointer: fine)')
		setIsHoverDevice(query.matches)
		const update = (event: MediaQueryListEvent) =>
			setIsHoverDevice(event.matches)
		query.addEventListener('change', update)

		return () => query.removeEventListener('change', update)
	}, [])

	const lifts = isHoverDevice && !shouldReduceMotion

	return (
		<section aria-labelledby={headingId} className={className}>
			<div className="mx-auto max-w-7xl px-4.5 py-24 sm:px-9 md:py-32">
				<div className="mx-auto mb-16 max-w-2xl text-center">
					<h2
						className="text-4xl font-normal tracking-[-0.04em] text-balance text-foreground md:text-5xl"
						id={headingId}
					>
						{title}
					</h2>
					<p className="mt-4 text-lg text-pretty text-muted-foreground">
						{lead}
					</p>
				</div>
				<ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
					{items.map((item, index) => (
						<motion.li
							className={cn(
								'relative rounded-xl border border-border bg-background p-6 transition-shadow focus-within:ring-2 focus-within:ring-ring',
								lifts && 'hover:shadow-md',
							)}
							initial={
								shouldReduceMotion
									? {
											opacity: 1,
										}
									: {
											opacity: 0,
											y: 20,
										}
							}
							key={item.id}
							transition={
								shouldReduceMotion
									? {
											duration: 0,
										}
									: {
											...SPRING,
											delay: index * STAGGER,
										}
							}
							viewport={{
								margin: '-100px',
								once: true,
							}}
							whileHover={
								lifts
									? {
											y: -4,
										}
									: undefined
							}
							whileInView={
								shouldReduceMotion
									? {
											opacity: 1,
										}
									: {
											opacity: 1,
											y: 0,
										}
							}
						>
							<div className="mb-4 inline-flex rounded-lg bg-foreground/5 p-2.5 text-foreground">
								{item.icon}
							</div>
							<h3 className="mb-2 font-semibold text-foreground">
								{item.title}
							</h3>
							<p className="text-sm leading-relaxed text-muted-foreground">
								{item.description}
							</p>
						</motion.li>
					))}
				</ul>
			</div>
		</section>
	)
}
