import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

export interface DialogTitleProps
	extends ComponentProps<typeof DialogPrimitive.Title> {}

// The dialog's heading, which also names the dialog.
export function DialogTitle(props: DialogTitleProps) {
	const { className, ...rest } = props

	return (
		<DialogPrimitive.Title
			className={cn('text-lg leading-none font-semibold', className)}
			data-slot="dialog-title"
			{...rest}
		/>
	)
}
