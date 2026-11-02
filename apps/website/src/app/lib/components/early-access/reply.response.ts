export function reply(
	status: number,
	code: string,
	retryAfter?: number,
	fieldErrors?: { email: string },
) {
	const headers: Record<string, string> = {
		'Cache-Control': 'no-store',
		'X-Content-Type-Options': 'nosniff',
	}
	if (retryAfter) headers['Retry-After'] = String(retryAfter)

	return Response.json(
		{
			ok: status === 202,
			code,
			...(fieldErrors ? { fieldErrors } : {}),
		},
		{
			status,
			headers,
		},
	)
}
