import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

// One section of a card, under a hairline rule.
export interface CardSectionProps {
	children: ReactNode
	className?: string
}

export function CardSection(props: CardSectionProps) {
	const { children, className } = props

	return (
		<div className={cn('mt-5 border-t border-border pt-4', className)}>
			{children}
		</div>
	)
}
