import type { JudgeUnavailableReason } from '@app/api/judge/judge.types'

// The platform's model gave no answer: its key is not set, or the provider
// could not be reached or refused the call. The judge answers 503, never a
// judgment.
export class JudgeModelUnavailable extends Error {
	readonly reason: JudgeUnavailableReason

	constructor(reason: JudgeUnavailableReason, options?: ErrorOptions) {
		super(`The judge model is unavailable: ${reason}`, options)
		this.name = 'JudgeModelUnavailable'
		this.reason = reason
	}
}
