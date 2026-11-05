const jsonContentType = /^application\/json(?:\s*;\s*charset=utf-8)?$/i

// A submission arrives as JSON in UTF-8, without a content encoding.
export function isPlainJSON(request: Request): boolean {
	if (request.headers.has('content-encoding')) return false

	return jsonContentType.test(request.headers.get('content-type') ?? '')
}
