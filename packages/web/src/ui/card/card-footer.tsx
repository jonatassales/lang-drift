import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

export interface CardFooterProps extends ComponentProps<'div'> {}

// The card's closing row, for its actions.
export function CardFooter(props: CardFooterProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn('flex items-center px-6 [.border-t]:pt-6', className)}
			data-slot="card-footer"
			{...rest}
		/>
	)
}
