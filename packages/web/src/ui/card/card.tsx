import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

export interface CardProps extends ComponentProps<'div'> {}

// shadcn/ui's new-york v4 card: a bordered surface its parts lay out.
export function Card(props: CardProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn(
				'flex flex-col gap-6 rounded-xl border border-border bg-card py-6 text-card-foreground shadow-sm',
				className,
			)}
			data-slot="card"
			{...rest}
		/>
	)
}
