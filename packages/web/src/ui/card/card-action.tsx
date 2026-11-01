import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

export interface CardActionProps extends ComponentProps<'div'> {}

// The control in the top right corner of the card's header.
export function CardAction(props: CardActionProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn(
				'col-start-2 row-span-2 row-start-1 self-start justify-self-end',
				className,
			)}
			data-slot="card-action"
			{...rest}
		/>
	)
}
