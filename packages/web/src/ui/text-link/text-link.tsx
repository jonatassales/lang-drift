import { Slot } from '@radix-ui/react-slot'
import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

// An inline link in the color of the text around it, underlined at a quarter
// of that color, with the site's focus outline.
export interface TextLinkProps extends ComponentProps<'a'> {
	/** Lends the link's look to its one child, a framework link or a button. */
	asChild?: boolean
}

export function TextLink(props: TextLinkProps) {
	const { asChild = false, className, ...rest } = props

	const Element = asChild ? Slot : 'a'

	return (
		<Element
			className={cn(
				'text-[13px] text-inherit underline decoration-current/25 underline-offset-5 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
				className,
			)}
			data-slot="text-link"
			{...rest}
		/>
	)
}
