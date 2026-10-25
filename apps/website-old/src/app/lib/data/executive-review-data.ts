export const targetIds = ['onboarding', 'integrations', 'permissions'] as const
export const evidenceKinds = [
	'spec',
	'decision',
	'verification',
	'feedback',
] as const
export const classificationFilters = [
	'Expected',
	'Unexpected',
	'unclassified',
] as const
export type TargetId = (typeof targetIds)[number]
export type EvidenceKind = (typeof evidenceKinds)[number]
export type Classification = 'Expected' | 'Unexpected'
export type ClassificationFilter = (typeof classificationFilters)[number]
export type Focus = 'all' | 'no-linked-evidence' | 'open-review'
export type QuestionId = 'overview' | Exclude<Focus, 'all'>
export interface ReviewFilters {
	targetId: TargetId | null
	classification: ClassificationFilter | null
	focus: Focus
}
export interface ReviewEvent {
	id: string
	targetId: TargetId
	classification: Classification | null
	week: 1 | 2 | 3 | 4
	sourceIds: readonly string[]
}
export interface ReviewSource {
	id: string
	eventId: string
	kind: EvidenceKind
}
export interface ReviewCheckpoint {
	id: string
	week: 1 | 2 | 3 | 4 | null
	sourceIds: readonly string[]
}
export interface ReviewItem {
	id: string
	eventId: string
	checkpoints: readonly ReviewCheckpoint[]
}
export interface ReviewTotals {
	events: number
	targets: number
	expected: number
	unexpected: number
	unclassified: number
	eventsWithEvidence: number
	sources: number
	openReviews: number
}

export const initialReviewFilters: ReviewFilters = {
	targetId: null,
	classification: null,
	focus: 'all',
}
export const questionPresets: Record<QuestionId, ReviewFilters> = {
	overview: initialReviewFilters,
	'no-linked-evidence': {
		targetId: null,
		classification: null,
		focus: 'no-linked-evidence',
	},
	'open-review': {
		targetId: null,
		classification: null,
		focus: 'open-review',
	},
}

export const reviewEvents: readonly ReviewEvent[] = [
	{
		id: 'E01',
		targetId: 'onboarding',
		classification: 'Expected',
		week: 1,
		sourceIds: ['E01-S', 'E01-D', 'E01-V', 'E01-F'],
	},
	{
		id: 'E02',
		targetId: 'onboarding',
		classification: 'Expected',
		week: 2,
		sourceIds: ['E02-S', 'E02-D', 'E02-V'],
	},
	{
		id: 'E03',
		targetId: 'onboarding',
		classification: 'Expected',
		week: 3,
		sourceIds: ['E03-S', 'E03-D', 'E03-F'],
	},
	{
		id: 'E04',
		targetId: 'onboarding',
		classification: 'Unexpected',
		week: 3,
		sourceIds: ['E04-S', 'E04-V', 'E04-F'],
	},
	{
		id: 'E05',
		targetId: 'onboarding',
		classification: 'Unexpected',
		week: 4,
		sourceIds: ['E05-S', 'E05-V'],
	},
	{
		id: 'E06',
		targetId: 'onboarding',
		classification: null,
		week: 4,
		sourceIds: [],
	},
	{
		id: 'E07',
		targetId: 'integrations',
		classification: 'Expected',
		week: 1,
		sourceIds: ['E07-S', 'E07-D', 'E07-V'],
	},
	{
		id: 'E08',
		targetId: 'integrations',
		classification: 'Expected',
		week: 2,
		sourceIds: ['E08-S', 'E08-D', 'E08-V', 'E08-F'],
	},
	{
		id: 'E09',
		targetId: 'integrations',
		classification: 'Expected',
		week: 3,
		sourceIds: ['E09-S', 'E09-D', 'E09-V'],
	},
	{
		id: 'E10',
		targetId: 'integrations',
		classification: 'Unexpected',
		week: 4,
		sourceIds: ['E10-S', 'E10-V', 'E10-F'],
	},
	{
		id: 'E11',
		targetId: 'permissions',
		classification: 'Expected',
		week: 2,
		sourceIds: ['E11-S', 'E11-D', 'E11-V', 'E11-F'],
	},
	{
		id: 'E12',
		targetId: 'permissions',
		classification: null,
		week: 4,
		sourceIds: [],
	},
]

// Evidence kind is explicit data, independent of the source ID suffix.
export const reviewSources: readonly ReviewSource[] = [
	{
		id: 'E01-S',
		eventId: 'E01',
		kind: 'spec',
	},
	{
		id: 'E01-D',
		eventId: 'E01',
		kind: 'decision',
	},
	{
		id: 'E01-V',
		eventId: 'E01',
		kind: 'verification',
	},
	{
		id: 'E01-F',
		eventId: 'E01',
		kind: 'feedback',
	},
	{
		id: 'E02-S',
		eventId: 'E02',
		kind: 'spec',
	},
	{
		id: 'E02-D',
		eventId: 'E02',
		kind: 'decision',
	},
	{
		id: 'E02-V',
		eventId: 'E02',
		kind: 'verification',
	},
	{
		id: 'E03-S',
		eventId: 'E03',
		kind: 'spec',
	},
	{
		id: 'E03-D',
		eventId: 'E03',
		kind: 'decision',
	},
	{
		id: 'E03-F',
		eventId: 'E03',
		kind: 'feedback',
	},
	{
		id: 'E04-S',
		eventId: 'E04',
		kind: 'spec',
	},
	{
		id: 'E04-V',
		eventId: 'E04',
		kind: 'verification',
	},
	{
		id: 'E04-F',
		eventId: 'E04',
		kind: 'feedback',
	},
	{
		id: 'E05-S',
		eventId: 'E05',
		kind: 'spec',
	},
	{
		id: 'E05-V',
		eventId: 'E05',
		kind: 'verification',
	},
	{
		id: 'E07-S',
		eventId: 'E07',
		kind: 'spec',
	},
	{
		id: 'E07-D',
		eventId: 'E07',
		kind: 'decision',
	},
	{
		id: 'E07-V',
		eventId: 'E07',
		kind: 'verification',
	},
	{
		id: 'E08-S',
		eventId: 'E08',
		kind: 'spec',
	},
	{
		id: 'E08-D',
		eventId: 'E08',
		kind: 'decision',
	},
	{
		id: 'E08-V',
		eventId: 'E08',
		kind: 'verification',
	},
	{
		id: 'E08-F',
		eventId: 'E08',
		kind: 'feedback',
	},
	{
		id: 'E09-S',
		eventId: 'E09',
		kind: 'spec',
	},
	{
		id: 'E09-D',
		eventId: 'E09',
		kind: 'decision',
	},
	{
		id: 'E09-V',
		eventId: 'E09',
		kind: 'verification',
	},
	{
		id: 'E10-S',
		eventId: 'E10',
		kind: 'spec',
	},
	{
		id: 'E10-V',
		eventId: 'E10',
		kind: 'verification',
	},
	{
		id: 'E10-F',
		eventId: 'E10',
		kind: 'feedback',
	},
	{
		id: 'E11-S',
		eventId: 'E11',
		kind: 'spec',
	},
	{
		id: 'E11-D',
		eventId: 'E11',
		kind: 'decision',
	},
	{
		id: 'E11-V',
		eventId: 'E11',
		kind: 'verification',
	},
	{
		id: 'E11-F',
		eventId: 'E11',
		kind: 'feedback',
	},
]

export const reviewItems: readonly ReviewItem[] = [
	{
		id: 'R01',
		eventId: 'E02',
		checkpoints: [
			{
				id: 'decision',
				week: 1,
				sourceIds: ['E02-D'],
			},
			{
				id: 'reviewed',
				week: 2,
				sourceIds: ['E02-V'],
			},
			{
				id: 'pending',
				week: null,
				sourceIds: [],
			},
		],
	},
	{
		id: 'R02',
		eventId: 'E04',
		checkpoints: [
			{
				id: 'scope',
				week: 1,
				sourceIds: ['E04-S'],
			},
			{
				id: 'observed',
				week: 3,
				sourceIds: ['E04-V', 'E04-F'],
			},
			{
				id: 'pending',
				week: null,
				sourceIds: [],
			},
		],
	},
	{
		id: 'R03',
		eventId: 'E06',
		checkpoints: [
			{
				id: 'inventory',
				week: 4,
				sourceIds: [],
			},
			{
				id: 'pending',
				week: null,
				sourceIds: [],
			},
		],
	},
	{
		id: 'R04',
		eventId: 'E12',
		checkpoints: [
			{
				id: 'inventory',
				week: 4,
				sourceIds: [],
			},
			{
				id: 'pending',
				week: null,
				sourceIds: [],
			},
		],
	},
]

export function getReviewsForEvents(
	events: readonly ReviewEvent[],
): ReviewItem[] {
	const ids = new Set(events.map((event) => event.id))
	return reviewItems.filter((review) => ids.has(review.eventId))
}

export function selectReviewEvents(
	filters: ReviewFilters,
	events = reviewEvents,
): ReviewEvent[] {
	const reviewedIds = new Set(reviewItems.map((review) => review.eventId))
	return events.filter((event) => {
		if (filters.targetId !== null && event.targetId !== filters.targetId)
			return false
		if (
			filters.classification !== null &&
			event.classification !==
				(filters.classification === 'unclassified'
					? null
					: filters.classification)
		)
			return false
		if (filters.focus === 'no-linked-evidence' && event.sourceIds.length !== 0)
			return false
		if (filters.focus === 'open-review' && !reviewedIds.has(event.id))
			return false
		return true
	})
}

export function getEventSources(
	event: ReviewEvent,
	kind?: EvidenceKind,
): ReviewSource[] {
	const ids = new Set(event.sourceIds)
	return reviewSources.filter(
		(source) =>
			ids.has(source.id) &&
			source.eventId === event.id &&
			(kind === undefined || source.kind === kind),
	)
}

export function aggregateReviewEvents(
	events: readonly ReviewEvent[],
): ReviewTotals {
	const unique = [...new Map(events.map((event) => [event.id, event])).values()]
	return {
		events: unique.length,
		targets: new Set(unique.map((event) => event.targetId)).size,
		expected: unique.filter((event) => event.classification === 'Expected')
			.length,
		unexpected: unique.filter((event) => event.classification === 'Unexpected')
			.length,
		unclassified: unique.filter((event) => event.classification === null)
			.length,
		eventsWithEvidence: unique.filter((event) => event.sourceIds.length > 0)
			.length,
		sources: new Set(unique.flatMap((event) => event.sourceIds)).size,
		openReviews: getReviewsForEvents(unique).length,
	}
}

export const reviewDistribution = targetIds.map((targetId) => {
	const totals = aggregateReviewEvents(
		reviewEvents.filter((event) => event.targetId === targetId),
	)
	return {
		targetId,
		total: totals.events,
		Expected: totals.expected,
		Unexpected: totals.unexpected,
		unclassified: totals.unclassified,
	}
})

export interface ExecutiveReviewMessages {
	eyebrow: string
	title: string
	description: string
	periodLabel: string
	questionLabel: string
	questions: Record<QuestionId, string>
	labels: {
		currentSelection: string
		target: string
		allTargets: string
		classification: string
		allClassifications: string
		reset: string
		relatedEvent: string
		nextDecision: string
		noLinkedSources: string
		eventReference: string
		reviewReference: string
		week: string
		pendingCondition: string
		earlyAccess: string
		selectionSummary: string
		resultStatus: string
	}
	metrics: {
		events: string
		targets: string
		openReviews: string
		eventsWithEvidence: string
		evidenceRatio: string
		noEvents: string
	}
	classifications: Record<ClassificationFilter, string>
	kinds: Record<EvidenceKind, string>
	chart: {
		title: string
		description: string
		scope: string
		segmentLabel: string
		targetLabel: string
	}
	matrix: {
		title: string
		description: string
		eventHeading: string
		linked: string
		notLinked: string
		linkedCellLabel: string
	}
	source: {
		title: string
		close: string
		titleTemplate: string
		referenceLabel: string
	}
	timeline: {
		title: string
		noReviews: string
		noReviewsLimit: string
		details: string
	}
	explanation: string
	answers: {
		overviewTitle: string
		overview: string
		noEvidenceTitle: string
		noEvidence: string
		openReviewTitle: string
		openReview: string
		manualSetup: string
		emptyTitle: string
		emptyDescription: string
	}
	targets: Record<TargetId, { label: string; description: string }>
	events: Record<string, { title: string; summary: string }>
	sources: Record<string, string>
	reviews: Record<
		string,
		{
			title: string
			summary: string
			nextDecision: string
			limit: string
			checkpoints: Record<string, { label: string; description: string }>
		}
	>
}
