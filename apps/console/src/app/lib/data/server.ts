import 'server-only'
import { headers } from 'next/headers'
import { cache } from 'react'

import { ConsoleError } from './api'
import { consoleHandler } from './console-handler'

/** React cache deduplicates reads only within this server render, never across users. */
export const readConsole = cache(async <T>(path: string): Promise<T> => {
	const incoming = await headers()

	const host = incoming.get('host') ?? ''
	const origin = new URL(`http://${host}`)
	if (!['localhost', '127.0.0.1', '[::1]'].includes(origin.hostname))
		throw new ConsoleError('local_only', 403)
	const response = await consoleHandler(
		new Request(new URL(`/api/console/${path}`, origin), { headers: incoming }),
	)

	const result = await response.json()

	if (!response.ok)
		throw new ConsoleError(
			result.error?.code ?? 'operation_failed',
			response.status,
		)

	return result as T
})
