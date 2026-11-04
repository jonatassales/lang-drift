import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

export interface CardContentProps extends ComponentProps<'div'> {}

// The card's body.
export function CardContent(props: CardContentProps) {
	const { className, ...rest } = props

	return (
		<div className={cn('px-6', className)} data-slot="card-content" {...rest} />
	)
}
