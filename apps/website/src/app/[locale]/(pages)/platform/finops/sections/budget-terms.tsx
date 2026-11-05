import {
	budgetTermIds,
	type FinopsMessages,
} from '@routes/platform/finops/content/finops-content.data'
import { DetailList } from 'web/ui/detail-list'

// The three words the funding model rests on, as definitions: the envelope,
// the allocation and the reservation.
interface BudgetTermsProps {
	copy: FinopsMessages['budget']['terms']
}

export function BudgetTerms(props: BudgetTermsProps) {
	const { copy } = props

	return (
		<DetailList
			className="mt-11 max-w-[960px]"
			items={budgetTermIds.map((id) => ({
				detail: (
					<span className="grid gap-1">
						<strong className="text-[15px] leading-snug font-medium">
							{copy[id].title}
						</strong>
						<span className="text-muted-foreground">
							{copy[id].description}
						</span>
					</span>
				),
				id,
				term: copy[id].term,
			}))}
		/>
	)
}
