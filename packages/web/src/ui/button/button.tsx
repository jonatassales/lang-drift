import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'
import type { ComponentProps } from 'react'

import { cn } from 'web/utilities'

// shadcn/ui's new-york v4 button, colored by the token roles. The tokens carry
// no white, so the destructive label takes the primary foreground, the white
// every filled control writes on.
const buttonVariants = cva(
	'inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*=size-])]:size-4',
	{
		variants: {
			variant: {
				default: 'bg-primary text-primary-foreground hover:bg-primary/90',
				destructive:
					'bg-destructive text-primary-foreground hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40',
				outline:
					'border border-border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50',
				secondary:
					'bg-secondary text-secondary-foreground hover:bg-secondary/80',
				ghost:
					'hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50',
				link: 'text-primary underline-offset-4 hover:underline',
				// The site's call to action: a filled pill that lifts on hover and takes
				// the site's focus outline instead of the shadcn ring.
				pill: 'gap-1.75 rounded-full bg-primary text-[13px] font-[650] text-primary-foreground transition-[translate,filter] duration-180 ease-[ease] hover:-translate-y-px hover:saturate-118 motion-reduce:transition-none motion-reduce:hover:translate-y-0 focus-visible:ring-0 focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
			},
			size: {
				default: 'h-9 px-4 py-2 has-[>svg]:px-3',
				xs: 'h-6 gap-1 rounded-md px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*=size-])]:size-3',
				sm: 'h-8 gap-1.5 rounded-md px-3 has-[>svg]:px-2.5',
				lg: 'h-10 rounded-md px-6 has-[>svg]:px-4',
				icon: 'size-9',
				'icon-xs': 'size-6 rounded-md [&_svg:not([class*=size-])]:size-3',
				'icon-sm': 'size-8',
				'icon-lg': 'size-10',
			},
		},
		compoundVariants: [
			{
				variant: 'pill',
				size: 'default',
				className: 'h-auto min-h-9.5 px-4.25 py-0 has-[>svg]:px-4.25',
			},
		],
		defaultVariants: {
			variant: 'default',
			size: 'default',
		},
	},
)

export interface ButtonProps
	extends ComponentProps<'button'>,
		VariantProps<typeof buttonVariants> {
	/** Lends the button's look to its one child, a link say. */
	asChild?: boolean
}

export function Button(props: ButtonProps) {
	const {
		asChild = false,
		className,
		size = 'default',
		variant = 'default',
		...rest
	} = props

	const Element = asChild ? Slot : 'button'

	return (
		<Element
			className={cn(
				buttonVariants({
					className,
					size,
					variant,
				}),
			)}
			data-size={size}
			data-slot="button"
			data-variant={variant}
			{...rest}
		/>
	)
}
