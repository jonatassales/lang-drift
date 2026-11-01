import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

export interface CardTitleProps extends ComponentProps<'div'> {}

// The card's title.
export function CardTitle(props: CardTitleProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn('leading-none font-semibold', className)}
			data-slot="card-title"
			{...rest}
		/>
	)
}
