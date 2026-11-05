import { POST } from '@api/book-a-demo/route'

import { demoRequest } from './demo-request.mock'

export function postDemoRequest(
	body: string,
	headers?: Record<string, string>,
) {
	return POST(demoRequest(body, headers))
}
