import { cn } from '@langdrift/web/utilities'

interface ShareBarProps {
	className?: string
	/** The filled share of the track, from 0 to 1. */
	share: number
}

// A two-pixel track whose fill grows to a share of it: a component's share
// of the CPU envelope, or the headroom the plan leaves. It stays mounted at
// zero before the share is known, so the growth is a transition.
export function ShareBar(props: ShareBarProps) {
	const { className, share } = props

	return (
		<span
			aria-hidden="true"
			className={cn(
				'block h-0.5 overflow-hidden rounded bg-foreground/5',
				className,
			)}
		>
			<span
				className="block h-full rounded bg-primary motion-safe:transition-[width] motion-safe:duration-700 motion-safe:ease-out"
				style={{ width: `${share * 100}%` }}
			/>
		</span>
	)
}
