import 'server-only'

import { extubRequest } from './extub'

/** Shared request boundary for route handlers and server-rendered read operations. */
export async function consoleHandler(request: Request) {
	const path = new URL(request.url).pathname.replace(/^\/api\/console\//u, '')

	return extubRequest(request, `console/${path}`)
}
