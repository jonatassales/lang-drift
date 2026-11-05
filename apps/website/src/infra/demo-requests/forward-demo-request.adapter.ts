import type { BookDemoSubmission } from '@routes/book-a-demo/form.types'

const timeoutMs = 3500

// Hands an accepted submission to the demo-requests service with the
// website's bearer token, and answers with the status the service gave.
// The address and the token are read here, where they are needed.
export async function forwardDemoRequest(
	submission: BookDemoSubmission,
	signal: AbortSignal,
): Promise<number> {
	const url = process.env.DEMO_REQUESTS_URL
	const token = process.env.DEMO_REQUESTS_TOKEN
	if (!url) throw new Error('Missing DEMO_REQUESTS_URL')
	if (!token) throw new Error('Missing DEMO_REQUESTS_TOKEN')

	const response = await fetch(`${url}/v1/demo-requests`, {
		body: JSON.stringify(submission),
		cache: 'no-store',
		headers: {
			authorization: `Bearer ${token}`,
			'content-type': 'application/json',
		},
		method: 'POST',
		redirect: 'error',
		signal: AbortSignal.any([signal, AbortSignal.timeout(timeoutMs)]),
	})

	const body = response.body

	if (body) await body.cancel()

	return response.status
}
