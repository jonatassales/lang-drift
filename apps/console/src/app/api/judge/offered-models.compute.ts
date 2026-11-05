// The models the console offers, from the comma-separated list its
// environment names, without blanks or repeats.
export function offeredModels(list: string | undefined): string[] {
	if (!list) return []

	const names = list
		.split(',')
		.map((name) => name.trim())
		.filter(Boolean)

	return [...new Set(names)]
}
