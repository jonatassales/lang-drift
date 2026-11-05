import { judgeReply } from '@app/api/judge/judge-reply.map'

// A question is a POST; the connection check is the route's GET.
export function interceptMethod(request: Request): Response | undefined {
	if (request.method === 'POST') return undefined

	return judgeReply(405, {
		code: 'METHOD_NOT_ALLOWED',
	})
}
