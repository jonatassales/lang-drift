import { answerSchemas } from './answer-schemas.data'
import { firstInvalidField } from './first-invalid-field.compute'
import { hasTokenCounts } from './has-token-counts.compute'
import type {
	JudgeRefusal,
	JudgeResponse,
	ModelPrices,
	ProviderResult,
	QuestionKind,
} from './judge.types'
import { judgeVersion } from './judge-version.data'
import { judgmentCost } from './judgment-cost.compute'

// What the route answers from what the provider returned: a judgment the
// platform stands behind, its text validated against the kind's schema and
// its usage priced at the list prices, or a refusal that names the cause and
// carries no text.
export function judgeResponse(
	result: ProviderResult,
	kind: QuestionKind,
	prices: ModelPrices,
): JudgeResponse | JudgeRefusal {
	const { stopReason, text, usage } = result

	if (!hasTokenCounts(usage)) {
		return {
			reason:
				'the provider reported a missing, negative or fractional token count',
		}
	}

	if (stopReason !== 'end_turn') {
		return {
			reason: `the model stopped with ${stopReason ?? 'no stop reason'}, not end_turn`,
		}
	}

	if (text.trim() === '') return { reason: 'the model answered with no text' }

	const tokensUsed = usage.inputTokens + usage.outputTokens

	if (tokensUsed === 0) {
		return {
			reason: 'the provider reported no tokens for a non-empty answer',
		}
	}

	let answer: unknown

	// The model's text is external input: text that does not parse is the
	// model's failure, answered as a refusal.
	try {
		answer = JSON.parse(text)
	} catch {
		return { reason: 'the answer is not JSON' }
	}

	const checked = answerSchemas[kind].safeParse(answer)

	if (!checked.success) {
		const field = firstInvalidField(checked.error, 'answer')

		return { reason: `the ${kind} answer breaks its schema at ${field}` }
	}

	return {
		cost: judgmentCost(usage, prices),
		judgeVersion,
		text,
		usage,
	}
}
