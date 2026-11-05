export function leaves(value: unknown, prefix = ''): Record<string, string> {
	if (typeof value === 'string') return { [prefix]: value }
	if (!value || typeof value !== 'object') return {}
	return Object.fromEntries(
		Object.entries(value).flatMap(([key, child]) =>
			Object.entries(leaves(child, prefix ? `${prefix}.${key}` : key)),
		),
	)
}
