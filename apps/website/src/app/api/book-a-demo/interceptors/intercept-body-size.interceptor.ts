import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { declaresOversizedBody } from '@components/early-access/declares-oversized-body.validation'

import { bodyLimit } from './body-limit.data'

// A declared length over the limit is refused before the body is read.
export function interceptBodySize(request: Request): Response | undefined {
	if (!declaresOversizedBody(request, bodyLimit)) return undefined

	return reply(413, 'INVALID_REQUEST')
}
