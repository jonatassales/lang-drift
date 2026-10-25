import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface CardSectionProps {
	children: ReactNode
	className?: string
}

// One section of the ticket card, under a hairline rule.
export function CardSection(props: CardSectionProps) {
	const { children, className } = props

	return (
		<div className={cn('mt-5 border-t border-border pt-4', className)}>
			{children}
		</div>
	)
}
