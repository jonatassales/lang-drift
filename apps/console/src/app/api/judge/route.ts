import { interceptBearer } from '@app/api/judge/interceptors/intercept-bearer.interceptor'
import { interceptBody } from '@app/api/judge/interceptors/intercept-body.interceptor'
import { interceptBodySize } from '@app/api/judge/interceptors/intercept-body-size.interceptor'
import { interceptContentType } from '@app/api/judge/interceptors/intercept-content-type.interceptor'
import { interceptMethod } from '@app/api/judge/interceptors/intercept-method.interceptor'
import { interceptModel } from '@app/api/judge/interceptors/intercept-model.interceptor'
import { interceptQuestion } from '@app/api/judge/interceptors/intercept-question.interceptor'
import { askJudgeModel } from '@infra/anthropic/ask-judge-model.adapter'
import { JudgeModelUnavailable } from '@infra/anthropic/judge-model-unavailable.type'

import type { ProviderResult } from './judge.types'
import { judgeReply } from './judge-reply.map'
import { judgeResponse } from './judge-response.compute'
import { judgeWireResponse } from './judge-wire-response.map'
import { logJudgment } from './log-judgment.adapter'

export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'
export const maxDuration = 120

// The connection check `dft agent connect` performs: 204 with the
// deployment's secret, 401 without it.
export function GET(request: Request): Response {
	const refused = interceptBearer(request)

	if (refused) return refused

	return new Response(null, {
		headers: {
			'Cache-Control': 'no-store',
		},
		status: 204,
	})
}

// A question passes one interceptor per concern, then the platform's model
// answers it with LangDrift's own key. The judge forwards only a judgment it
// validated, or a refusal, with what the call used and cost.
export async function POST(request: Request): Promise<Response> {
	const refused =
		interceptMethod(request) ??
		interceptBearer(request) ??
		interceptContentType(request) ??
		interceptBodySize(request)

	if (refused) return refused

	const read = await interceptBody(request)

	if (read instanceof Response) return read

	const accepted = interceptQuestion(read.body)

	if (accepted instanceof Response) return accepted

	const { question } = accepted
	const offered = interceptModel(question.model)

	if (offered instanceof Response) return offered

	let result: ProviderResult

	try {
		result = await askJudgeModel(question, request.signal)
	} catch (error) {
		if (!(error instanceof JudgeModelUnavailable)) throw error

		console.error(`Judge model unavailable: ${error.reason}`)

		return judgeReply(503, {
			code: 'UNAVAILABLE',
			reason: error.reason,
		})
	}

	const outcome = judgeResponse(result, question.kind, offered.prices)
	const answer = judgeWireResponse(outcome, result.usage, offered.prices)

	logJudgment({
		cents: answer.cost.cents,
		inputTokens: answer.usage.inputTokens,
		kind: question.kind,
		model: question.model,
		outputTokens: answer.usage.outputTokens,
		promptVersion: question.promptVersion,
		refused: answer.refused,
	})

	return judgeReply(200, answer)
}
