import { consoleRequest } from './api'

/** Action-created promises are consumed by use() within their local Suspense boundary. */
export type ReadResult<T> =
	| { data: T; error?: never }
	| { data?: never; error: unknown }
export async function readResult<T>(
	path: string,
	signal?: AbortSignal,
): Promise<ReadResult<T>> {
	try {
		const data = await consoleRequest<T>(path, { signal })

		return { data }
	} catch (error) {
		return { error }
	}
}
