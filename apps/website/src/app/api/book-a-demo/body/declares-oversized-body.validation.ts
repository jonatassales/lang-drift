// A declared length is refused before the body is read, and one that is not a
// decimal count of bytes counts as oversized. The read itself stays bounded,
// because a declaration can lie.
export function declaresOversizedBody(
	request: Request,
	limit: number,
): boolean {
	const declared = request.headers.get('content-length')

	if (!declared) return false

	return !/^\d+$/.test(declared) || Number(declared) > limit
}
