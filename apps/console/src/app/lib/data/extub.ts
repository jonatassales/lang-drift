import { consoleHostAllowed, extubBase } from './extub-link'

/** Server-only bridge: stub mode never forwards browser credentials or invokes native workloads. */
export async function extubRequest(
	request: Request,
	path: string,
): Promise<Response> {
	try {
		const url = new URL(request.url)
		const host = request.headers.get('host') ?? url.host
		const browser = new URL(`${url.protocol}//${host}`)
		const origin = request.headers.get('origin')
		const token = process.env.EXTUB_TOKEN?.trim()
		const crossSite =
			Boolean(origin && origin !== browser.origin) ||
			request.headers.get('sec-fetch-site') === 'cross-site'
		if (!consoleHostAllowed(browser.hostname, token) || crossSite)
			return Response.json({ error: { code: 'local_only' } }, { status: 403 })
		const base = extubBase(process.env.EXTUB_URL, token)
		if (!base)
			return Response.json({ error: { code: 'local_only' } }, { status: 503 })
		let body: Uint8Array | undefined
		if (!['GET', 'HEAD'].includes(request.method) && request.body) {
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
				error: { code: 'extub_unavailable' },
			},
			{ status: 503 },
		)
	}
}
