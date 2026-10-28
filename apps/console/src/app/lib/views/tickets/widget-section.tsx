import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface WidgetSectionProps {
	children: ReactNode
	className?: string
}

// One section of a ticket widget, under a hairline rule.
export function WidgetSection(props: WidgetSectionProps) {
	const { children, className } = props

	return (
		<div className={cn('mt-4 border-t border-border pt-3', className)}>
			{children}
		</div>
	)
}
