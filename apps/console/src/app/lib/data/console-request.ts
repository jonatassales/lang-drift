import { ConsoleError } from './console-error'
import { consolePathSafe } from './entity-id'

export async function consoleRequest<T>(
	path: string,
	init?: RequestInit,
): Promise<T> {
	// An id outside the alphabet names no record, so nothing is requested.
	if (!consolePathSafe(path)) throw new ConsoleError('not_found', 404)
	try {
		const response = await fetch(`/api/console/${path}`, {
			credentials: 'same-origin',
			cache: 'no-store',
			...init,
			headers: {
				'Content-Type': 'application/json',
				...init?.headers,
			},
		})

		const result: unknown = await response.json()

		if (!response.ok) {
			const code =
				typeof result === 'object' &&
				result !== null &&
				'error' in result &&
				typeof result.error === 'object' &&
				result.error !== null &&
				'code' in result.error &&
				typeof result.error.code === 'string'
					? result.error.code
					: 'operation_failed'
			throw new ConsoleError(code, response.status)
		}

		return result as T
	} catch (error) {
		if (
			error instanceof ConsoleError ||
			(error instanceof DOMException && error.name === 'AbortError')
		)
			throw error
		throw new ConsoleError('connection_failed', 0)
	}
}
