import 'server-only'

import { interceptBody } from '@api/book-a-demo/interceptors/intercept-body.interceptor'
import { interceptBodySize } from '@api/book-a-demo/interceptors/intercept-body-size.interceptor'
import { interceptContentType } from '@api/book-a-demo/interceptors/intercept-content-type.interceptor'
import { interceptMethod } from '@api/book-a-demo/interceptors/intercept-method.interceptor'
import { interceptOrigin } from '@api/book-a-demo/interceptors/intercept-origin.interceptor'
import { interceptSchema } from '@api/book-a-demo/interceptors/intercept-schema.interceptor'
import { forwardDemoRequest } from '@infra/demo-requests/forward-demo-request.adapter'

import { reply } from './book-demo-reply.response'

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
