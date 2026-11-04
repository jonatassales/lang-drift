import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

export interface EyebrowProps {
	as?: 'span' | 'dt'
	children: ReactNode
	className?: string
}

// A small mono label over a value or a card: the term of a figure, the kind
// of a record, the period of a report.
export function Eyebrow(props: EyebrowProps) {
	const { as: Element = 'span', children, className } = props

	return (
		<Element
			className={cn(
				'font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase',
				className,
			)}
		>
			{children}
		</Element>
	)
}
