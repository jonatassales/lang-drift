export function demoRequest(
	body: string,
	headers: Record<string, string> = {},
	method = 'POST',
) {
	return new Request('http://localhost:3000/api/book-a-demo/', {
		body,
		headers: {
			'content-type': 'application/json',
			...headers,
		},
		method,
	})
}
