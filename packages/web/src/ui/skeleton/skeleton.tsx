import { cn } from '@langdrift/web/utilities'
import type { ComponentProps } from 'react'

export interface SkeletonProps extends ComponentProps<'div'> {}

// shadcn/ui's new-york v4 skeleton: a pulsing block that holds the place of
// content still loading. It is decorative, so assistive technology skips it,
// and it rests still for visitors who ask for reduced motion.
export function Skeleton(props: SkeletonProps) {
	const { className, ...rest } = props

	return (
		<div
			aria-hidden
			className={cn(
				'animate-pulse rounded-md bg-accent motion-reduce:animate-none',
				className,
			)}
			data-slot="skeleton"
			{...rest}
		/>
	)
}
