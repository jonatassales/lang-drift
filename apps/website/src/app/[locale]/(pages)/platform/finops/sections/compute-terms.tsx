import { computeTermIds } from '@routes/platform/finops/content/finops-content.data'
import type { FinopsMessages } from '@routes/platform/finops/content/finops-content.types'
import { DetailList } from 'web/ui/detail-list'

// The three words hardware spend is governed by, as definitions: the compute
// envelope, the quota and cleanup.
interface ComputeTermsProps {
	copy: FinopsMessages['budget']['terms']
}

export function ComputeTerms(props: ComputeTermsProps) {
	const { copy } = props

	return (
		<DetailList
			className="mt-11 max-w-[960px]"
			items={computeTermIds.map((id) => ({
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
