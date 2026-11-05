import { planPhaseIds } from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { Skeleton } from 'web/ui/skeleton'
import { cn } from 'web/utilities'

// The card stack's place held while the terminal types the command: the
// ticket card open over the three collapsed headers, pulsing until the
// stack fades in over it. It is decorative, so assistive technology skips
// it; the stack beneath stays hidden until the command is entered.
interface SessionCardsSkeletonProps {
	/** True until the command is entered and the stack fades in. */
	pending: boolean
}

export function SessionCardsSkeleton(props: SessionCardsSkeletonProps) {
	const { pending } = props

	return (
		<div
			aria-hidden="true"
			className={cn(
				'pointer-events-none absolute inset-0 flex flex-col gap-3 motion-safe:transition-opacity motion-safe:duration-500',
				pending ? 'visible opacity-100' : 'invisible opacity-0',
			)}
		>
			<div className="flex min-h-0 flex-1 flex-col gap-5 overflow-hidden rounded-2xl border border-border bg-card px-5 py-4 sm:px-6">
				<div className="flex items-center justify-between gap-4">
					<Skeleton className="h-3 w-16" />
					<Skeleton className="h-3 w-24" />
				</div>
				<Skeleton className="h-5 w-3/4" />
				<div className="grid gap-3">
					<Skeleton className="h-3 w-full" />
					<Skeleton className="h-3 w-5/6" />
					<Skeleton className="h-3 w-2/3" />
				</div>
			</div>
			{planPhaseIds.slice(1).map((phase) => (
				<div
					className="flex min-h-11 shrink-0 items-center justify-between gap-4 rounded-2xl border border-border bg-card px-5 py-3 sm:px-6"
					key={phase}
				>
					<Skeleton className="h-3 w-20" />
					<Skeleton className="h-3 w-28" />
				</div>
			))}
		</div>
	)
}
