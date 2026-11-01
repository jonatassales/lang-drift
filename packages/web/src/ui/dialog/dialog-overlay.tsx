import { cn } from '@langdrift/web/utilities'
import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

export interface DialogOverlayProps
	extends ComponentProps<typeof DialogPrimitive.Overlay> {}

// The layer behind a modal dialog. The tokens carry no black, so it veils the
// page with the background role instead.
export function DialogOverlay(props: DialogOverlayProps) {
	const { className, ...rest } = props

	return (
		<DialogPrimitive.Overlay
			className={cn('fixed inset-0 z-50 bg-background/80', className)}
			data-slot="dialog-overlay"
			{...rest}
		/>
	)
}
