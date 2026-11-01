import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

export interface DialogHeaderProps extends ComponentProps<'div'> {}

// The dialog's title and description, stacked.
export function DialogHeader(props: DialogHeaderProps) {
	const { className, ...rest } = props

	return (
		<div
			className={cn('flex flex-col gap-2 text-center sm:text-left', className)}
			data-slot="dialog-header"
			{...rest}
		/>
	)
}
