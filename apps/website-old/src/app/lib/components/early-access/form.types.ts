import type { Draft } from 'immer'

export type EmailErrorCode = 'emailRequired' | 'emailInvalid' | 'emailTooLong'

export interface EarlyAccessFormState {
	email: string
	hasInteractedWithEmail: boolean
	fieldErrors: { email: EmailErrorCode | null }
	retryAfter?: number
	status:
		| 'idle'
		| 'verifying'
		| 'submitting'
		| 'success'
		| 'rate-limited'
		| 'error'
		| 'challenge-error'
}

export type StateUpdater<State> = (
	recipe: (draft: Draft<State>) => void,
) => void

export type EarlyAccessActionResult =
	| { ok: true }
	| { ok: false; code: 'VALIDATION_ERROR'; fieldErrors?: { email?: string } }
	| { ok: false; code: 'CHALLENGE_FAILED' }
	| { ok: false; code: 'UNAVAILABLE' }
	| { ok: false; code: 'RATE_LIMITED'; retryAfter: number }
