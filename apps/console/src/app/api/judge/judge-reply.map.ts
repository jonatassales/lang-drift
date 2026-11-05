import type { JudgeErrorBody, JudgeWireResponse } from './judge.types'

// Every answer of the route is JSON that no cache keeps.
export function judgeReply(
	status: number,
	body: JudgeErrorBody | JudgeWireResponse,
): Response {
	return Response.json(body, {
		headers: {
			'Cache-Control': 'no-store',
			'X-Content-Type-Options': 'nosniff',
		},
		status,
	})
}
