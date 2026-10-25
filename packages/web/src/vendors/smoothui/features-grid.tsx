'use client'

import { GlowHoverCards } from '@langdrift/web/vendors/smoothui/glow-hover-card'
import { motion, useReducedMotion } from 'motion/react'
import type { ReactNode } from 'react'

// Adapted from SmoothUI's Features Grid block (smoothui.dev): a heading and a
// grid of cards that rise into view and glow under the pointer through
// SmoothUI's Glow Hover Card. The content comes from the caller.

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

const HIDDEN = {
	opacity: 0,
	y: 20,
}

const SHOWN = {
	opacity: 1,
	y: 0,
}

const VIEWPORT = {
	margin: '-100px',
	once: true,
}

export function FeaturesGrid(props: FeaturesGridProps) {
	const { className, headingId, items, lead, title } = props

	const shouldReduceMotion = useReducedMotion()

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
				<motion.div
					initial={shouldReduceMotion ? false : HIDDEN}
					transition={SPRING}
					viewport={VIEWPORT}
					whileInView={SHOWN}
				>
					<GlowHoverCards
						cardClassName="relative rounded-xl border border-border bg-background p-6 focus-within:ring-2 focus-within:ring-ring"
						className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4"
						items={items.map((item) => ({
							content: (
								<>
									<div className="mb-4 inline-flex rounded-lg bg-foreground/5 p-2.5 text-foreground">
										{item.icon}
									</div>
									<h3 className="mb-2 font-semibold text-foreground">
										{item.title}
									</h3>
									<p className="text-sm leading-relaxed text-muted-foreground">
										{item.description}
									</p>
								</>
							),
							id: item.id,
						}))}
					/>
				</motion.div>
			</div>
		</section>
	)
}
