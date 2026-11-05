import { declaresOversizedBody } from '@api/book-a-demo/body/declares-oversized-body.validation'
import { reply } from '@api/book-a-demo/book-demo-reply.response'

import { bodyLimit } from './body-limit.data'

// A declared length over the limit is refused before the body is read.
export function interceptBodySize(request: Request): Response | undefined {
	if (!declaresOversizedBody(request, bodyLimit)) return undefined

	return reply(413, 'INVALID_REQUEST')
}
