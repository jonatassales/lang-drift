const jsonContentType = /^application\/json(?:\s*;\s*charset=utf-8)?$/i

// A submission arrives as JSON in UTF-8, without a content encoding.
export function isPlainJSON(request: Request): boolean {
	if (request.headers.has('content-encoding')) return false

	return jsonContentType.test(request.headers.get('content-type') ?? '')
}

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

// JSON.parse silently replaces duplicate properties. Tokenize first to reject
// duplicate decoded keys, including escaped aliases, at every nesting level.
export function parseSubmissionJSON(text: string): unknown {
	const parsed: unknown = JSON.parse(text)
	const tokens = text.match(/"(?:[^"\\]|\\.)*"|[{}[\],:]|[^\s{}[\],:]+/g) ?? []
	const stack: Array<{ keys: Set<string>; key: boolean; object: boolean }> = []
	for (const token of tokens) {
		if (token === '{' || token === '[') {
			if (stack.length >= 16) throw new Error('invalid_request')
			stack.push({
				keys: new Set(),
				key: token === '{',
				object: token === '{',
			})
		} else if (token === '}' || token === ']') {
			stack.pop()
		} else {
			const current = stack.at(-1)
			if (!current?.object) continue
			if (token === ',') current.key = true
			else if (token === ':') current.key = false
			else if (current.key) {
				const key: string = JSON.parse(token)
				if (current.keys.has(key)) throw new Error('invalid_request')
				current.keys.add(key)
				current.key = false
			}
		}
	}
	if (parsed === null || Array.isArray(parsed) || typeof parsed !== 'object') {
		throw new Error('invalid_request')
	}

	return parsed
}
