export function toInitials(name: string): string {
	if (!name) return ''

	return name
		.trim()
		.split(/\s+/)
		.slice(0, 2)
		.map((word) => word[0])
		.join('')
		.toUpperCase()
}
