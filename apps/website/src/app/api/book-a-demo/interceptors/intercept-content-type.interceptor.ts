import { isPlainJSON } from '@app/api/book-a-demo/body/is-plain-json.validation'
import { reply } from '@components/book-a-demo/book-demo-reply.response'

// A demo request is plain JSON: no encoding, no other media type.
export function interceptContentType(request: Request): Response | undefined {
	if (isPlainJSON(request)) return undefined

	return reply(415, 'INVALID_REQUEST')
}
