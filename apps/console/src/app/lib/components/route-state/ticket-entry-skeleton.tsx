import { Skeleton } from 'web/ui/skeleton'

// A ticket entry while it loads: the card below lg, the row from lg.
export function TicketEntrySkeleton() {
	return (
		<>
			<div className="grid grid-cols-1 gap-2 rounded-xl border border-border p-4 lg:hidden">
				<div className="flex items-center justify-between gap-3">
					<div className="flex gap-2">
						<Skeleton className="h-4 w-12" />
						<Skeleton className="h-4 w-24" />
					</div>
					<Skeleton className="h-6 w-24 rounded-full" />
				</div>
				<Skeleton className="h-4 w-4/5" />
				<Skeleton className="h-3 w-2/5" />
				<Skeleton className="mt-3 h-3 w-full" />
			</div>
			<div className="hidden items-center gap-4 rounded-xl border border-border px-4 py-3 lg:flex">
				<Skeleton className="h-4 w-12 shrink-0" />
				<div className="grid min-w-0 flex-1 gap-1.5">
					<Skeleton className="h-4 w-3/5" />
					<Skeleton className="h-3 w-2/5" />
				</div>
				<div className="grid w-36 shrink-0 justify-items-end gap-1.5">
					<Skeleton className="h-6 w-24 rounded-full" />
					<Skeleton className="h-3 w-20" />
				</div>
			</div>
		</>
	)
}
