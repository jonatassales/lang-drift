import {
	chainPillarIds,
	type FinopsMessages,
	failureScenarioIds,
	failureStatusTones,
} from '@routes/platform/finops/content/finops-content'
import { DetailList } from 'web/ui/detail-list'

interface FailureScenariosProps {
	copy: FinopsMessages['failures']
}

// The two ways delivery leaves the plan, side by side: each walks the chain
// from intent to outcome and ends in the state it leaves behind.
export function FailureScenarios(props: FailureScenariosProps) {
	const { copy } = props

	return (
		<div className="mt-11 grid gap-3 lg:grid-cols-2">
			{failureScenarioIds.map((id) => {
				const scenario = copy.scenarios[id]

				return (
					<article
						className="flex flex-col gap-3 rounded-2xl border border-border bg-card p-6 sm:p-8"
						key={id}
					>
						<h3 className="m-0 text-[22px] leading-tight font-medium tracking-[-0.02em]">
							{scenario.title}
						</h3>
						<p className="m-0 text-[15px] leading-[1.6] text-foreground">
							{scenario.context}
						</p>
						<DetailList
							className="mt-2 border-t border-border pt-1"
							items={[
								...chainPillarIds.map((pillar) => ({
									detail: scenario[pillar],
									id: pillar,
									term: copy.pillars[pillar],
								})),
								{
									detail: scenario.status.detail,
									id: 'status',
									term: scenario.status.term,
									tone: failureStatusTones[id],
								},
							]}
						/>
					</article>
				)
			})}
		</div>
	)
}
