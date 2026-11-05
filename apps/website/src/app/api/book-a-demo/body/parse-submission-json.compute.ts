// JSON.parse silently replaces duplicate properties. Tokenize first to reject
// duplicate decoded keys, including escaped aliases, at every nesting level.
interface JsonFrame {
	expectsKey: boolean
	isObject: boolean
	keys: Set<string>
}

const maxDepth = 16

export function parseSubmissionJSON(text: string): unknown {
	const parsed: unknown = JSON.parse(text)
	const tokens = text.match(/"(?:[^"\\]|\\.)*"|[{}[\],:]|[^\s{}[\],:]+/g) ?? []
	const stack: JsonFrame[] = []

	for (const token of tokens) {
		const opens = token === '{' || token === '['
		const closes = token === '}' || token === ']'

		if (opens) {
			if (stack.length >= maxDepth) throw new Error('invalid_request')

			stack.push({
				expectsKey: token === '{',
				isObject: token === '{',
				keys: new Set(),
			})

			continue
		}

		if (closes) {
			stack.pop()

			continue
		}

		const current = stack.at(-1)

		if (!current?.isObject) continue

		if (token === ',') {
			current.expectsKey = true

			continue
		}

		if (token === ':') {
			current.expectsKey = false

			continue
		}

		if (!current.expectsKey) continue

		const key: string = JSON.parse(token)

		if (current.keys.has(key)) throw new Error('invalid_request')

		current.keys.add(key)
		current.expectsKey = false
	}

	const notAnObject =
		parsed === null || Array.isArray(parsed) || typeof parsed !== 'object'

	if (notAnObject) {
		throw new Error('invalid_request')
	}

	return parsed
}
