export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

/** Authenticated project secrets; credentials and queries remain server-owned. */
export {
	consoleHandler as GET,
	consoleHandler as POST,
} from '@data/console-handler'
