import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { isPlainJSON } from '@components/early-access/is-plain-json.validation'

// A demo request is plain JSON: no encoding, no other media type.
export function interceptContentType(request: Request): Response | undefined {
	if (isPlainJSON(request)) return undefined

	return reply(415, 'INVALID_REQUEST')
}
