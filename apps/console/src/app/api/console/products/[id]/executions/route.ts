import { consoleHandler } from '@data/console-handler'
export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

/** Authenticated project executions; credentials and queries remain server-owned. */
export async function GET(request: Request) {
	return consoleHandler(request)
}
