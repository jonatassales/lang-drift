// Stable local content structure; these IDs are not an API or classification model.
export const visionLoopScenarios = [
	{
		id: 'priority-shift',
		loop: ['request', 'change', 'review'],
		evidence: ['priority-target', 'priority-request', 'priority-flow'],
	},
	{
		id: 'product-trade-off',
		loop: ['decision', 'delivery', 'review-condition'],
		evidence: ['tradeoff-target', 'tradeoff-decision', 'tradeoff-release'],
	},
	{
		id: 'new-opportunity',
		loop: ['discovery', 'experiment', 'proposed-direction'],
		evidence: [
			'opportunity-interviews',
			'opportunity-prototype',
			'opportunity-proposal',
		],
	},
] as const

export type ScenarioId = (typeof visionLoopScenarios)[number]['id']

export interface EvidenceCopy {
	kindLabel: string
	title: string
	summary: string
	sourceContext: string
}

export interface ScenarioCopy {
	tabLabel: string
	context: string
	vision: { target: string; audience: string; successCriterion: string }
	loop: Record<string, { label: string; description: string }>
	evidence: Record<string, EvidenceCopy>
	interpretation: { title: string; description: string; openQuestion: string }
}

export interface VisionLoopMessages {
	heading: string
	introduction: string
	selectorLabel: string
	labels: {
		vision: string
		loop: string
		evidence: string
		target: string
		audience: string
		successCriterion: string
		supportingRecords: string
		sourceContext: string
		interpretation: string
		openQuestion: string
	}
	driftExplanation: string
	continueToAttribution: string
	scenarios: {
		[Definition in (typeof visionLoopScenarios)[number] as Definition['id']]: ScenarioCopy & {
			loop: Record<
				Definition['loop'][number],
				{ label: string; description: string }
			>
			evidence: Record<Definition['evidence'][number], EvidenceCopy>
		}
	}
}
