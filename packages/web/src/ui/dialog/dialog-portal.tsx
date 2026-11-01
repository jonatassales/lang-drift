import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogPortalProps
	extends ComponentProps<typeof DialogPrimitive.Portal> {}

// Renders the overlay and the content at the end of the document's body.
export function DialogPortal(props: DialogPortalProps) {
	const { children, ...rest } = props

	return (
		<DialogPrimitive.Portal data-slot="dialog-portal" {...rest}>
			{children}
		</DialogPrimitive.Portal>
	)
}
