import {
	chainPillarIds,
	failureScenarioIds,
} from '@routes/platform/finops/content/finops-content.data'
import type { FinopsMessages } from '@routes/platform/finops/content/finops-content.types'
import { PillarScenarios } from '@routes/platform/product-page/pillar-scenarios'

// The two ways hardware leaves the plan, as tabs over the chain from intent
// to outcome.
interface FailureScenariosProps {
	copy: FinopsMessages['failures']
}

export function FailureScenarios(props: FailureScenariosProps) {
	const { copy } = props

	const pillars = chainPillarIds.map((id) => ({
		id,
		label: copy.pillars[id],
	}))
	const scenarios = failureScenarioIds.map((id) => {
		const scenario = copy.scenarios[id]

		return {
			cells: {
				intent: scenario.intent,
				delivery: scenario.delivery,
				evidence: scenario.evidence,
				outcome: scenario.outcome,
			},
			context: scenario.context,
			id,
			tabLabel: scenario.tabLabel,
		}
	})

	return (
		<PillarScenarios
			pillars={pillars}
			scenarios={scenarios}
			selectorLabel={copy.selector}
		/>
	)
}
