import { firstInvalidField } from '@app/api/judge/first-invalid-field.compute'
import type { JudgeQuestion } from '@app/api/judge/judge.types'
import { judgeReply } from '@app/api/judge/judge-reply.map'
import { judgeRequestSchema } from '@app/api/judge/judge-request.validation'

export interface AcceptedQuestion {
	question: JudgeQuestion
}

// The contract's request schema decides, and a refusal names the field.
export function interceptQuestion(body: unknown): Response | AcceptedQuestion {
	const result = judgeRequestSchema.safeParse(body)

	if (!result.success) {
		return judgeReply(400, {
			code: 'INVALID_QUESTION',
			field: firstInvalidField(result.error, 'question'),
		})
	}

	return {
		question: result.data,
	}
}
