import 'server-only'

import { consolePathSafe } from './entity-id'
import { extubRequest } from './extub'

/** Shared request boundary for route handlers and server-rendered read operations. */
export async function consoleHandler(request: Request) {
	const path = new URL(request.url).pathname.replace(/^\/api\/console\//u, '')
	if (!consolePathSafe(path))
		return Response.json({ error: { code: 'not_found' } }, { status: 404 })

	return extubRequest(request, `console/${path}`)
}
