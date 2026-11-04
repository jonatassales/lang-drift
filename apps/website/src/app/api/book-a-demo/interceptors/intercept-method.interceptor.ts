import { reply } from '@components/book-a-demo/book-demo-reply.response'

// Only a POST carries a demo request.
export function interceptMethod(request: Request): Response | undefined {
	if (request.method === 'POST') return undefined

	return reply(405, 'INVALID_REQUEST')
}
