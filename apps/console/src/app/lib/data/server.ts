import 'server-only'
import { headers } from 'next/headers'
import { cache } from 'react'

import { ConsoleError } from './api'
import { consoleHandler } from './console-handler'
import { consoleHostAllowed, serverRead } from './console-stubs-link'

/** React cache deduplicates reads only within this server render, never across users. */
export const readConsole = cache(
	async <Result>(path: string): Promise<Result> => {
		const incoming = await headers()

		const host = incoming.get('host') ?? ''
		const origin = new URL(`http://${host}`)
		if (!consoleHostAllowed(origin.hostname, process.env.CONSOLE_STUBS_TOKEN))
			throw new ConsoleError('local_only', 403)

		const response = await consoleHandler(serverRead(host, path))

		const result = await response.json()

		if (!response.ok)
			throw new ConsoleError(
				result.error?.code ?? 'operation_failed',
				response.status,
			)

		return result as Result
	},
)
