import { consoleRequest } from './console-request'

export function mutation<Result>(
	path: string,
	method: 'POST' | 'PATCH' | 'DELETE',
	body?: unknown,
): Promise<Result> {
	return consoleRequest<Result>(path, {
		method,
		body: JSON.stringify(body ?? {}),
	})
}
