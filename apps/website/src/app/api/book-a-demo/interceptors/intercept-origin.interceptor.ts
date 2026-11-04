import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { isForeignRequest } from '@components/book-a-demo/is-foreign-request.validation'

// Only the site's own pages post demo requests.
export function interceptOrigin(request: Request): Response | undefined {
	if (!isForeignRequest(request)) return undefined

	return reply(403, 'INVALID_ORIGIN')
}
