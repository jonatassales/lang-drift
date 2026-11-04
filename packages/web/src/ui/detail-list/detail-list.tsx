import type { ReactNode } from 'react'

import { Eyebrow } from 'web/ui/eyebrow'
import { cn } from 'web/utilities'

export type DetailTone = 'default' | 'success' | 'warning' | 'danger'

export interface DetailListItem {
	detail: ReactNode
	id: string
	term: string
	tone?: DetailTone
}

export interface DetailListProps {
	className?: string
	items: readonly DetailListItem[]
}

const tones: Record<DetailTone, string> = {
	default: 'text-foreground',
	success: 'text-status-success',
	warning: 'text-status-warning',
	danger: 'text-status-danger',
}

// A definition list of term and detail rows, for receipts, reports and
// summaries. A tone colors the detail with a status role.
export function DetailList(props: DetailListProps) {
	const { className, items } = props

	return (
		<dl className={cn('m-0 grid', className)}>
			{items.map((item) => (
				<div
					className="grid gap-1 border-t border-border py-3 first:border-t-0 sm:grid-cols-[minmax(0,0.4fr)_minmax(0,1fr)] sm:gap-4"
					key={item.id}
				>
					<Eyebrow as="dt">{item.term}</Eyebrow>
					<dd
						className={cn(
							'm-0 text-[13px] leading-[1.55]',
							tones[item.tone ?? 'default'],
						)}
					>
						{item.detail}
					</dd>
				</div>
			))}
		</dl>
	)
}
