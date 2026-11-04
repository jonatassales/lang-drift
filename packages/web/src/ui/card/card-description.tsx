import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

export interface CardDescriptionProps extends ComponentProps<'div'> {}

// The muted line under the card's title.
export function CardDescription(props: CardDescriptionProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn('text-sm text-muted-foreground', className)}
			data-slot="card-description"
			{...rest}
		/>
	)
}
