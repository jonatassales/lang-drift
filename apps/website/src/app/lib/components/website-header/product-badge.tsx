import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface ProductBadgeProps {
	children: ReactNode
	className?: string
}

export function ProductBadge(props: ProductBadgeProps) {
	const { children, className } = props

	return (
		<span
			className={cn(
				'inline-flex items-center rounded-full border border-border bg-muted/70 px-2 py-0.5 text-[10px] font-medium leading-none text-muted-foreground',
				className,
			)}
		>
			{children}
		</span>
	)
}
