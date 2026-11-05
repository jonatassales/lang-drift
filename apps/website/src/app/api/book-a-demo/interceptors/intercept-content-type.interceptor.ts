import { isPlainJSON } from '@api/book-a-demo/body/is-plain-json.compute'
import { bookDemoReply } from '@api/book-a-demo/book-demo-reply.map'

// A demo request is plain JSON: no encoding, no other media type.
export function interceptContentType(request: Request): Response | undefined {
	if (isPlainJSON(request)) return undefined

	return bookDemoReply(415, 'INVALID_REQUEST')
}
