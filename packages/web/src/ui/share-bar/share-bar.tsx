import { cn } from '@langdrift/web/utilities'

export interface ShareBarProps {
	className?: string
	/** The filled share of the track, from 0 to 1; more than 1 fills it. */
	share: number
	/** A share that breaks a limit takes the warning role. */
	tone?: 'default' | 'warning'
}

// A two-pixel track whose fill grows to a share of it: a component's share
// of the CPU envelope, the headroom the plan leaves, or a reservation's share
// of the budget. It stays mounted at zero before the share is known, so the
// growth is a transition.
export function ShareBar(props: ShareBarProps) {
	const { className, share, tone = 'default' } = props

	return (
		<span
			aria-hidden="true"
			className={cn(
				'block h-0.5 overflow-hidden rounded bg-foreground/5',
				className,
			)}
		>
			<span
				className={cn(
					'block h-full rounded motion-safe:transition-[width] motion-safe:duration-700 motion-safe:ease-out',
					tone === 'warning' ? 'bg-status-warning' : 'bg-primary',
				)}
				style={{ width: `${Math.min(share, 1) * 100}%` }}
			/>
		</span>
	)
}
