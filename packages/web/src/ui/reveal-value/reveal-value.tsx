import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

import type { RevealValueLabels } from './reveal-value-labels.type'
import type { RevealValueState } from './reveal-value-state.type'

// One value as it is revealed: a dash while nothing is known, a pulsing
// skeleton while it resolves, then the value fading in. A value not known
// is never drawn as zero. Screen readers get the state as words.
export interface RevealValueProps {
	children?: ReactNode
	className?: string
	labels: RevealValueLabels
	state: RevealValueState
	/** The skeleton's width class while the value resolves. */
	width?: string
}

export function RevealValue(props: RevealValueProps) {
	const { children, className, labels, state, width = 'w-24' } = props

	if (state === 'unset') {
		return (
			<span className={className}>
				<span aria-hidden="true">—</span>
				<span className="sr-only">{labels.notSet}</span>
			</span>
		)
	}

	if (state === 'resolving') {
		return (
			<span className={cn('block', className)}>
				<span
					aria-hidden="true"
					className={cn(
						'block h-3 rounded bg-foreground/10 animate-pulse motion-reduce:animate-none',
						width,
					)}
				/>
				<span className="sr-only">{labels.resolving}</span>
			</span>
		)
	}

	return (
		<span
			className={cn(
				'motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0',
				className,
			)}
		>
			{children}
		</span>
	)
}
