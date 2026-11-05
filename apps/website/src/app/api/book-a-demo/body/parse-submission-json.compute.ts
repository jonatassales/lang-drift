// JSON.parse silently replaces duplicate properties. Tokenize first to reject
// duplicate decoded keys, including escaped aliases, at every nesting level.
interface JsonFrame {
	keys: Set<string>
	key: boolean
	object: boolean
}

export function parseSubmissionJSON(text: string): unknown {
	const parsed: unknown = JSON.parse(text)
	const tokens = text.match(/"(?:[^"\\]|\\.)*"|[{}[\],:]|[^\s{}[\],:]+/g) ?? []
	const stack: JsonFrame[] = []

	for (const token of tokens) {
		const opens = token === '{' || token === '['
		const closes = token === '}' || token === ']'

		if (opens) {
			if (stack.length >= 16) throw new Error('invalid_request')

			stack.push({
				keys: new Set(),
				key: token === '{',
				object: token === '{',
			})
			continue
		}

		if (closes) {
			stack.pop()
			continue
		}

		const current = stack.at(-1)

		if (!current?.object) continue

		if (token === ',') {
			current.key = true
			continue
		}

		if (token === ':') {
			current.key = false
			continue
		}

		if (!current.key) continue

		const key: string = JSON.parse(token)

		if (current.keys.has(key)) throw new Error('invalid_request')

		current.keys.add(key)
		current.key = false
	}

	const notAnObject =
		parsed === null || Array.isArray(parsed) || typeof parsed !== 'object'

	if (notAnObject) {
		throw new Error('invalid_request')
	}

	return parsed
}
