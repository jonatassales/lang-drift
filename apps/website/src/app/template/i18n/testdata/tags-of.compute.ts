export function tagsOf(message: string): string[] {
	return Array.from(
		new Set(
			Array.from(message.matchAll(/<([A-Za-z]\w*)>/g), (match) => match[1]),
		),
	).sort()
}
