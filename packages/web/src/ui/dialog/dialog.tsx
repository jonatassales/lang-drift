import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogProps
	extends ComponentProps<typeof DialogPrimitive.Root> {}

// shadcn/ui's new-york v4 dialog, on Radix: the root that holds whether the
// dialog is open and whether it is modal.
export function Dialog(props: DialogProps) {
	const { children, ...rest } = props

	return (
		<DialogPrimitive.Root data-slot="dialog" {...rest}>
			{children}
		</DialogPrimitive.Root>
	)
}
