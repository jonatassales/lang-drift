import 'server-only'

import { interceptBody } from '@app/api/book-a-demo/interceptors/intercept-body.interceptor'
import { interceptBodySize } from '@app/api/book-a-demo/interceptors/intercept-body-size.interceptor'
import { interceptContentType } from '@app/api/book-a-demo/interceptors/intercept-content-type.interceptor'
import { interceptMethod } from '@app/api/book-a-demo/interceptors/intercept-method.interceptor'
import { interceptOrigin } from '@app/api/book-a-demo/interceptors/intercept-origin.interceptor'
import { interceptSchema } from '@app/api/book-a-demo/interceptors/intercept-schema.interceptor'
import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { forwardDemoRequest } from '@infra/demo-requests/forward-demo-request.adapter'

export const runtime = 'nodejs'

// A demo request passes the edge's interceptors, one per concern, and only
// then reaches the demo-requests service. While that service is unreachable
// or unconfigured, the visitor learns that requests are paused.
export async function POST(request: Request): Promise<Response> {
	const refused =
		interceptMethod(request) ??
		interceptOrigin(request) ??
		interceptContentType(request) ??
		interceptBodySize(request)
	if (refused) return refused

	const read = await interceptBody(request)

	if (read instanceof Response) return read

	const accepted = interceptSchema(read.body)
	if (accepted instanceof Response) return accepted

	let status: number
	try {
		status = await forwardDemoRequest(accepted.submission, request.signal)
	} catch (error) {
		const reason = error instanceof Error ? error.message : String(error)
		console.error(`Demo request not forwarded: ${reason}`)

		return reply(503, 'UNAVAILABLE')
	}

	if (status !== 202) {
		console.error(`Demo request refused upstream with status ${status}`)

		return reply(503, 'UNAVAILABLE')
	}

	return reply(200, 'OK')
}
