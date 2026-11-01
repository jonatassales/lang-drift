import { cn } from '@langdrift/web/utilities'
import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogDescriptionProps
	extends ComponentProps<typeof DialogPrimitive.Description> {}

// The muted line under the title, which also describes the dialog.
export function DialogDescription(props: DialogDescriptionProps) {
	const { className, ...rest } = props

	return (
		<DialogPrimitive.Description
			className={cn('text-sm text-muted-foreground', className)}
			data-slot="dialog-description"
			{...rest}
		/>
	)
}
