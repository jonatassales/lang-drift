import type { z } from 'zod'

import type { judgeRequestSchema } from './judge-request.validation'

// One question as the contract carries it: the kind, the planner model, the
// system text and the prompt the CLI versions, and the token bound.
export type JudgeQuestion = z.output<typeof judgeRequestSchema>

// The four questions the planner of the LangDrift CLI asks, one kind each.
export type QuestionKind = JudgeQuestion['kind']

// What a call consumed, in tokens, as the provider reported it.
export interface ProviderUsage {
	inputTokens: number
	outputTokens: number
}

// What the provider answered: the text of its text blocks joined, what the
// call consumed and why the model stopped.
export interface ProviderResult {
	stopReason: string | null
	text: string
	usage: ProviderUsage
}

// A model's list prices, in cents per million tokens.
export interface ModelPrices {
	inputCentsPerMillion: number
	outputCentsPerMillion: number
}

// What a judgment cost and who says so.
export interface JudgeCost {
	cents: number
	illustrative: true
	provenance: 'console'
}

// A judgment the platform stands behind: the answer text validated against
// the kind's schema, the usage, its cost and the judge that answered.
export interface JudgeResponse {
	cost: JudgeCost
	judgeVersion: string
	text: string
	usage: ProviderUsage
}

// A judgment the platform refuses to forward, with the cause.
export interface JudgeRefusal {
	reason: string
}

// The contract's response on the wire, a judgment or a refusal.
export interface JudgeWireResponse extends JudgeResponse {
	reason: string
	refused: boolean
}

export type JudgeErrorCode =
	| 'INVALID_QUESTION'
	| 'METHOD_NOT_ALLOWED'
	| 'MODEL_NOT_OFFERED'
	| 'PAYLOAD_TOO_LARGE'
	| 'UNAUTHORIZED'
	| 'UNAVAILABLE'
	| 'UNSUPPORTED_MEDIA_TYPE'

export type JudgeUnavailableReason =
	| 'judge-not-configured'
	| 'platform-key-missing'
	| 'provider-unavailable'

// A refused request: the code, and the field, the offered models or the
// cause when the code needs one.
export interface JudgeErrorBody {
	code: JudgeErrorCode
	field?: string
	offered?: string[]
	reason?: JudgeUnavailableReason
}
