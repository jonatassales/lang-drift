import { consoleRequest } from './api'

/** Action-created promises are consumed by use() within their local Suspense boundary. */
export type ReadResult<Result> =
	| { data: Result; error?: never }
	| { data?: never; error: unknown }

export async function readResult<Result>(
	path: string,
	signal?: AbortSignal,
): Promise<ReadResult<Result>> {
	try {
		const data = await consoleRequest<Result>(path, { signal })

		return { data }
	} catch (error) {
		return { error }
	}
}
