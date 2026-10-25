import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

import type { MicroruntimesMessages } from '@data/microruntimes-content'

export type CardValueState = 'unset' | 'resolving' | 'set'

interface CardValueProps {
	children?: ReactNode
	className?: string
	labels: Pick<MicroruntimesMessages['session']['card'], 'notSet' | 'resolving'>
	state: CardValueState
	/** The skeleton's width class while the value resolves. */
	width?: string
}

// One value on the ticket card as the session reveals it: a dash while
// nothing is set, a pulsing skeleton while the console resolves it, then the
// value. Screen readers get the state as words.
export function CardValue(props: CardValueProps) {
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

	return <span className={className}>{children}</span>
}
