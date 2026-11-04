import 'server-only'

import {
	consoleHostAllowed,
	consoleStubsBase,
	crossSite,
} from './console-stubs-link'

/** Server-only bridge: stub mode never forwards browser credentials or invokes native workloads. */
export async function consoleStubsRequest(
	request: Request,
	path: string,
): Promise<Response> {
	try {
		const url = new URL(request.url)
		const host = request.headers.get('host') ?? url.host
		const browser = new URL(`${url.protocol}//${host}`)
		const token = process.env.CONSOLE_STUBS_TOKEN?.trim()
		const hostAllowed = consoleHostAllowed(browser.hostname, token)
		const isCrossSite = crossSite(request)
		if (!hostAllowed || isCrossSite)
			return Response.json({ error: { code: 'local_only' } }, { status: 403 })

		const consoleStubsUrl = process.env.CONSOLE_STUBS_URL
		if (!consoleStubsUrl)
			return Response.json({ error: { code: 'local_only' } }, { status: 503 })

		const base = consoleStubsBase(consoleStubsUrl, token)
		if (!base)
			return Response.json({ error: { code: 'local_only' } }, { status: 503 })

		let body: Uint8Array | undefined
		const hasPayload =
			!['GET', 'HEAD'].includes(request.method) && Boolean(request.body)
		if (hasPayload && request.body) {
			const reader = request.body.getReader()
			const chunks: Uint8Array[] = []
			let total = 0
			for (;;) {
				const next = await reader.read()

				if (next.done) break
				total += next.value.length
				if (total > 65536) {
					await reader.cancel()

					return Response.json(
						{ error: { code: 'request_too_large' } },
						{ status: 413 },
					)
				}
				chunks.push(next.value)
			}
			body = new Uint8Array(total)
			let offset = 0
			for (const chunk of chunks) {
				body.set(chunk, offset)
				offset += chunk.length
			}
		}
		const target = new URL(`/${path}`, base)
		target.search = url.search
		const headers: Record<string, string> = body
			? { 'content-type': 'application/json' }
			: {}
		if (token) headers.authorization = `Bearer ${token}`

		const response = await fetch(target, {
			method: request.method,
			headers,
			body: body as BodyInit | undefined,
			cache: 'no-store',
			redirect: 'error',
			signal: AbortSignal.any([request.signal, AbortSignal.timeout(15000)]),
		})

		return new Response(response.body, {
			status: response.status,
			headers: {
				'content-type': 'application/json',
				'cache-control': 'no-store',
			},
		})
	} catch {
		return Response.json(
			{
				source: 'stub',
				error: { code: 'console_stubs_unavailable' },
			},
			{ status: 503 },
		)
	}
}
