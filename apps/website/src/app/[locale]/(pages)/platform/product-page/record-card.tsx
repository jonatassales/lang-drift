import { DetailList, type DetailListItem } from 'web/ui/detail-list'
import { Eyebrow } from 'web/ui/eyebrow'

// A titled record on a product page, such as a weekly report or an execution
// ledger: its title, the period it covers when it has one, and its rows.
interface RecordCardProps {
	items: readonly DetailListItem[]
	period?: string
	title: string
}

export function RecordCard(props: RecordCardProps) {
	const { items, period, title } = props

	return (
		<article className="mt-11 rounded-2xl border border-border bg-card p-6 sm:p-8">
			<div className="flex flex-wrap items-baseline justify-between gap-x-6 gap-y-2 border-b border-border pb-4">
				<strong className="text-lg font-medium tracking-[-0.02em]">
					{title}
				</strong>
				{period && <Eyebrow className="text-[10px]">{period}</Eyebrow>}
			</div>
			<DetailList className="mt-1" items={items} />
		</article>
	)
}
