import { bookDemoReply } from '@api/book-a-demo/book-demo-reply.response'
import { isForeignRequest } from '@api/book-a-demo/is-foreign-request.validation'

// Only the site's own pages post demo requests.
export function interceptOrigin(request: Request): Response | undefined {
	if (!isForeignRequest(request)) return undefined

	return bookDemoReply(403, 'INVALID_ORIGIN')
}
