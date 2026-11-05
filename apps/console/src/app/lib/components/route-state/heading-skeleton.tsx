'use client'
import { Skeleton } from 'web/ui/skeleton'

export function HeadingSkeleton() {
	return (
		<div className="mb-6 grid grid-cols-1 gap-3">
			<Skeleton className="h-7 w-56 max-w-full" />
			<Skeleton className="h-4 w-96 max-w-full" />
		</div>
	)
}
