import { consoleRequest } from './console-request'

export function mutation<T>(
	path: string,
	method: 'POST' | 'PATCH' | 'DELETE',
	body?: unknown,
): Promise<T> {
	return consoleRequest<T>(path, {
		method,
		body: JSON.stringify(body ?? {}),
	})
}
