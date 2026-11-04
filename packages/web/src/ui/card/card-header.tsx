import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

export interface CardHeaderProps extends ComponentProps<'div'> {}

// The card's heading row: its title and description, and an action beside them.
export function CardHeader(props: CardHeaderProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn(
				'@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6',
				className,
			)}
			data-slot="card-header"
			{...rest}
		/>
	)
}
