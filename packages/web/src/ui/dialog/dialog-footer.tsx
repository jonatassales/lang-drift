import { Button } from '@langdrift/web/ui/button'
import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

import { DialogClose } from './dialog-close'

export interface DialogFooterProps extends ComponentProps<'div'> {
	/** Labels an outline close button after the actions; without it there is none. */
	closeLabel?: string
}

// The dialog's actions, stacked on small screens and in a row from sm up.
export function DialogFooter(props: DialogFooterProps) {
	const { children, className, closeLabel, ...rest } = props

	return (
		<div
			className={cn(
				'flex flex-col-reverse gap-2 sm:flex-row sm:justify-end',
				className,
			)}
			data-slot="dialog-footer"
			{...rest}
		>
			{children}
			{closeLabel ? (
				<DialogClose asChild>
					<Button variant="outline">{closeLabel}</Button>
				</DialogClose>
			) : null}
		</div>
	)
}
