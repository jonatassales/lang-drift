export function argumentsOf(message: string): string[] {
	return Array.from(
		new Set(
			Array.from(
				message.matchAll(/\{([A-Za-z]\w*)\s*(?:,|\})/g),
				(match) => match[1],
			),
		),
	).sort()
}
