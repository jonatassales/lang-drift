import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface CardProps {
	children: ReactNode
	className?: string
}

export function Card(props: CardProps) {
	const { children, className } = props

	return (
		<section
			className={cn('rounded-[10px] border border-border bg-card', className)}
		>
			{children}
		</section>
	)
}
