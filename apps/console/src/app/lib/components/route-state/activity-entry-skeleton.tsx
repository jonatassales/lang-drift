import { Skeleton } from '@ui/skeleton'

const figures = ['first', 'second', 'third', 'fourth']

// An Executions or History entry while it loads: its runtime and state over
// the line of ids and the labelled figures.
export function ActivityEntrySkeleton() {
	return (
		<div className="grid grid-cols-1 gap-3 rounded-xl border border-border p-4">
			<div className="flex items-center justify-between gap-3">
				<Skeleton className="h-4 w-40 max-w-full" />
				<Skeleton className="h-6 w-20 rounded-full" />
			</div>
			<Skeleton className="h-3 w-3/5" />
			<div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
				{figures.map((key) => (
					<Skeleton key={key} className="h-8 w-full" />
				))}
			</div>
		</div>
	)
}
