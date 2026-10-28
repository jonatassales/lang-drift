import { cn } from '@langdrift/web/utilities'

interface SkeletonProps {
	className?: string
}

export function Skeleton(props: SkeletonProps) {
	const { className } = props

	return (
		<div
			aria-hidden
			className={cn(
				'animate-pulse rounded-md bg-foreground/10 motion-reduce:animate-none',
				className,
			)}
		/>
	)
}
