export function shape(value: unknown): string {
	if (Array.isArray(value))
		return `[${value.map((item) => shape(item)).join(',')}]`
	if (value !== null && typeof value === 'object') {
		return `{${Object.keys(value)
			.map((key) => `${key}:${shape(Reflect.get(value, key))}`)
			.join(',')}}`
	}

	return typeof value
}
