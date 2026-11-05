import { bookDemoReply } from '@api/book-a-demo/book-demo-reply.response'

// Only a POST carries a demo request.
export function interceptMethod(request: Request): Response | undefined {
	if (request.method === 'POST') return undefined

	return bookDemoReply(405, 'INVALID_REQUEST')
}
