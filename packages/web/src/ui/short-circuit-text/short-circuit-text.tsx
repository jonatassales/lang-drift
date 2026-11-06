'use client'

import { type ReactNode, useRef } from 'react'

import { cn } from 'web/utilities'

import { useDischarge } from './use-discharge.hook'

const sides = ['left', 'right'] as const

// Text in a short circuit: it hums in place, and at irregular moments the
// current arcs out of both sides in jagged bolts, a white core in a blue glow,
// while the word jerks and flickers. `useDischarge` owns the motion.
export interface ShortCircuitTextProps {
	children: ReactNode
}

export function ShortCircuitText(props: ShortCircuitTextProps) {
	const { children } = props

	const wordRef = useRef<HTMLSpanElement>(null)
	const discharge = useDischarge(wordRef)

	return (
		<span className="relative inline-block">
			<span className="inline-block" ref={wordRef}>
				{children}
			</span>
			{sides.map((side) => (
				<svg
					aria-hidden="true"
					className={cn(
						'pointer-events-none absolute top-1/2 h-[1.4em] w-[1.2em] -translate-y-1/2 overflow-visible drop-shadow-md drop-shadow-rainbow-5',
						side === 'left' ? 'right-full' : 'left-full -scale-x-100',
					)}
					key={side}
					preserveAspectRatio="none"
					viewBox="0 0 100 100"
				>
					{discharge[side].map((path) => (
						<g fill="none" key={path} strokeLinejoin="bevel">
							<path
								className="stroke-rainbow-5 opacity-50"
								d={path}
								strokeWidth={5}
								vectorEffect="non-scaling-stroke"
							/>
							<path
								className="stroke-primary-foreground"
								d={path}
								strokeWidth={1.5}
								vectorEffect="non-scaling-stroke"
							/>
						</g>
					))}
				</svg>
			))}
		</span>
	)
}
