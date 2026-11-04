import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

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
