import type {
	categoryIds,
	evidenceIds,
	flowStepIds,
	limitIds,
	principleIds,
	safeSeatFaqIds,
	sourceIds,
	whyIds,
} from './safe-seat-content.data'

// The shape of the safe-seat catalog, keyed by the page's ids.

export type WhyId = (typeof whyIds)[number]
export type EvidenceId = (typeof evidenceIds)[number]
export type CategoryId = (typeof categoryIds)[number]
export type FlowStepId = (typeof flowStepIds)[number]
export type PrincipleId = (typeof principleIds)[number]
export type LimitId = (typeof limitIds)[number]
export type SourceId = (typeof sourceIds)[number]
export type SafeSeatFaqId = (typeof safeSeatFaqIds)[number]

export interface SafeSeatTitledCopy {
	title: string
	description: string
}

export interface SafeSeatStepCopy extends SafeSeatTitledCopy {
	eyebrow: string
}

export interface SafeSeatTermCopy {
	term: string
	detail: string
}

export interface SafeSeatEvidenceItemCopy extends SafeSeatTitledCopy {
	source: string
	limit: string
}

export interface SafeSeatQuestionCopy {
	question: string
	answer: string
}

export interface SafeSeatSectionCopy {
	kicker: string
	title: string
	lead: string
}

export interface SafeSeatHeroCopy {
	comingSoon: string
	taglineStart: string
	taglineEnd: string
	lead: string
	note: string
}

export interface SafeSeatWhyCopy extends SafeSeatSectionCopy {
	items: Record<WhyId, SafeSeatTitledCopy>
}

export interface SafeSeatEvidenceCopy extends SafeSeatSectionCopy {
	limitLabel: string
	sourceLink: string
	followUpLink: string
	items: Record<EvidenceId, SafeSeatEvidenceItemCopy>
}

export interface SafeSeatReservedCopy extends SafeSeatSectionCopy {
	categories: Record<CategoryId, SafeSeatTitledCopy>
	never: string
}

export interface SafeSeatFlowCopy extends SafeSeatSectionCopy {
	steps: Record<FlowStepId, SafeSeatStepCopy>
}

export interface SafeSeatPrinciplesCopy extends SafeSeatSectionCopy {
	items: Record<PrincipleId, SafeSeatTermCopy>
}

export interface SafeSeatLimitsCopy extends SafeSeatSectionCopy {
	items: Record<LimitId, SafeSeatTermCopy>
}

export interface SafeSeatSourcesCopy extends SafeSeatSectionCopy {
	items: Record<SourceId, string>
}

export interface SafeSeatFaqCopy {
	items: Record<SafeSeatFaqId, SafeSeatQuestionCopy>
}

export interface SafeSeatMessages {
	skip: string
	hero: SafeSeatHeroCopy
	why: SafeSeatWhyCopy
	evidence: SafeSeatEvidenceCopy
	reserved: SafeSeatReservedCopy
	flow: SafeSeatFlowCopy
	principles: SafeSeatPrinciplesCopy
	limits: SafeSeatLimitsCopy
	sources: SafeSeatSourcesCopy
	faq: SafeSeatFaqCopy
}
