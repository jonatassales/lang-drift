import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogTriggerProps
	extends ComponentProps<typeof DialogPrimitive.Trigger> {}

// The control that opens the dialog and takes focus back when it closes.
export function DialogTrigger(props: DialogTriggerProps) {
	const { children, ...rest } = props

	return (
		<DialogPrimitive.Trigger data-slot="dialog-trigger" {...rest}>
			{children}
		</DialogPrimitive.Trigger>
	)
}
