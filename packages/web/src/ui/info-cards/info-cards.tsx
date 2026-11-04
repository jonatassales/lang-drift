import type { ReactNode } from 'react'

import { Eyebrow } from 'web/ui/eyebrow'
import { cn } from 'web/utilities'

export interface InfoCard {
	description: string
	eyebrow?: string
	footer?: ReactNode
	id: string
	title: string
}

export interface InfoCardsProps {
	className?: string
	columns?: 3 | 4
	items: readonly InfoCard[]
	ordered?: boolean
}

const columnClasses: Record<NonNullable<InfoCardsProps['columns']>, string> = {
	3: 'sm:grid-cols-2 lg:grid-cols-3',
	4: 'sm:grid-cols-2 lg:grid-cols-4',
}

// A grid of cards, each with an optional eyebrow, a title, a
// description and an optional footer the caller fills (a link, a note).
// Ordered when the cards are steps.
export function InfoCards(props: InfoCardsProps) {
	const { className, columns = 3, items, ordered = false } = props

	const List = ordered ? 'ol' : 'ul'

	return (
		<List
			className={cn(
				'm-0 grid list-none gap-3 p-0',
				columnClasses[columns],
				className,
			)}
		>
			{items.map((item) => (
				<li
					className="flex flex-col gap-2 rounded-2xl border border-border bg-card p-6"
					key={item.id}
				>
					{item.eyebrow ? <Eyebrow>{item.eyebrow}</Eyebrow> : null}
					<h3 className="m-0 text-lg leading-snug font-medium tracking-[-0.02em] text-foreground">
						{item.title}
					</h3>
					<p className="m-0 text-[13px] leading-[1.55] text-muted-foreground">
						{item.description}
					</p>
					{item.footer ? (
						<div className="mt-auto pt-3">{item.footer}</div>
					) : null}
				</li>
			))}
		</List>
	)
}
