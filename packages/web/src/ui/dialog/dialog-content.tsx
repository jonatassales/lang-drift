import { cn } from '@langdrift/web/utilities'
import { X } from '@phosphor-icons/react/ssr'
import * as DialogPrimitive from '@radix-ui/react-dialog'
import type { ComponentProps } from 'react'

import { DialogClose } from './dialog-close'
import { DialogOverlay } from './dialog-overlay'
import { DialogPortal } from './dialog-portal'

export interface DialogContentProps
	extends ComponentProps<typeof DialogPrimitive.Content> {
	/** Names the close button in the corner; without it there is none. */
	closeLabel?: string
}

// The dialog itself, centered over its overlay. The corner close button needs
// a label from the app's catalog, so it shows only when one is given.
export function DialogContent(props: DialogContentProps) {
	const { children, className, closeLabel, ...rest } = props

	return (
		<DialogPortal>
			<DialogOverlay />
			<DialogPrimitive.Content
				className={cn(
					'fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border border-border bg-background p-6 shadow-lg outline-none sm:max-w-lg',
					className,
				)}
				data-slot="dialog-content"
				{...rest}
			>
				{children}
				{closeLabel ? (
					<DialogClose className="absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*=size-])]:size-4">
						<X aria-hidden />
						<span className="sr-only">{closeLabel}</span>
					</DialogClose>
				) : null}
			</DialogPrimitive.Content>
		</DialogPortal>
	)
}
