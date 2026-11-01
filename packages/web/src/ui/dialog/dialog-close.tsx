import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogCloseProps
	extends ComponentProps<typeof DialogPrimitive.Close> {}

// A control that closes the dialog.
export function DialogClose(props: DialogCloseProps) {
	const { children, ...rest } = props

	return (
		<DialogPrimitive.Close data-slot="dialog-close" {...rest}>
			{children}
		</DialogPrimitive.Close>
	)
}
