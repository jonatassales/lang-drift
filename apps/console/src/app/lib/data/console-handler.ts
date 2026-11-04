import 'server-only'

import { consoleStubsRequest } from './console-stubs'
import { consolePathSafe } from './entity-id'

/** Shared request boundary for route handlers and server-rendered read operations. */
export async function consoleHandler(request: Request) {
	const path = new URL(request.url).pathname.replace(/^\/api\/console\//u, '')
	if (!consolePathSafe(path))
		return Response.json({ error: { code: 'not_found' } }, { status: 404 })

	return consoleStubsRequest(request, `console/${path}`)
}
