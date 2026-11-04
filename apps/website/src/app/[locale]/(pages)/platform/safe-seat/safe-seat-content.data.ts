import type { EvidenceId, SourceId } from './safe-seat-messages.types'

// Safe Seat page content: the ids and links the page renders. Sentences and
// citations live in the safe-seat catalog; this module holds what a reader
// does not read as prose. Safe Seat is in development, so the page describes
// the planned design.

export const whyIds = ['thinking', 'irony', 'delegate'] as const

export const evidenceIds = [
	'bainbridge',
	'casner',
	'shen-tamkin',
	'metr',
	'lee',
	'feng',
] as const

export const categoryIds = ['core-domain', 'hard-part', 'new-stack'] as const

export const flowStepIds = [
	'plan',
	'boundary',
	'dev-ticket',
	'attestation',
	'delivery',
] as const

export const principleIds = [
	'private',
	'team',
	'no-ranking',
	'floor',
	'overrides',
] as const

export const limitIds = [
	'beyond',
	'long-term',
	'share',
	'no-call',
	'governance',
	'analogy',
] as const

export const sourceIds = [
	'bainbridge',
	'casner',
	'shen-tamkin',
	'metr-2025',
	'metr-2026',
	'lee',
	'feng',
	'safo-13002',
	'safo-17007',
] as const

export const safeSeatFaqIds = [
	'quota',
	'manager',
	'ai',
	'why-five',
	'available',
	'outside',
] as const

// Each source's primary record: a DOI, the paper or the agency's own PDF.
export const sourceLinks: Record<SourceId, string> = {
	bainbridge: 'https://doi.org/10.1016/0005-1098(83)90046-8',
	casner: 'https://doi.org/10.1177/0018720814535628',
	'shen-tamkin': 'https://arxiv.org/abs/2601.20245',
	'metr-2025': 'https://arxiv.org/abs/2507.09089',
	'metr-2026': 'https://metr.org/blog/2026-02-24-uplift-update/',
	lee: 'https://doi.org/10.1145/3706598.3713778',
	feng: 'https://arxiv.org/abs/2510.07435',
	'safo-13002':
		'https://www.faa.gov/sites/faa.gov/files/other_visit/aviation_industry/airline_operators/airline_safety/SAFO13002.pdf',
	'safo-17007': 'https://www.faa.gov/sites/faa.gov/files/2022-11/SAFO17007.pdf',
}

// The source each evidence card links to.
export const evidenceSources: Record<EvidenceId, SourceId> = {
	bainbridge: 'bainbridge',
	casner: 'casner',
	'shen-tamkin': 'shen-tamkin',
	metr: 'metr-2025',
	lee: 'lee',
	feng: 'feng',
}

// A later source a card also reports, linked beside its primary one.
export const evidenceFollowUps: Partial<Record<EvidenceId, SourceId>> = {
	metr: 'metr-2026',
}
