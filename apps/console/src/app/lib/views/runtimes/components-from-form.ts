import type { Component } from '@data/types'

export function componentsFromForm(
	form: FormData,
	candidates: Component[],
): Component[] | null {
	const selected = new Set(form.getAll('selectedComponents').map(String))
	if (!selected.size) return null

	const filtered = candidates.filter((candidate) =>
		selected.has(candidate.path),
	)
	const result: Component[] = []

	for (const candidate of filtered) {
		const prefix = `component:${candidate.path}:`
		let parsed: unknown
		try {
			parsed = JSON.parse(String(form.get(`${prefix}args`)))
		} catch {
			return null
		}

		if (
			!Array.isArray(parsed) ||
			!parsed.every((argument) => typeof argument === 'string')
		) {
			return null
		}

		result.push({
			...candidate,
			image: String(form.get(`${prefix}image`)),
			port: Number(form.get(`${prefix}port`)),
			readinessPath: String(form.get(`${prefix}readinessPath`)),
			build: {
				executable: String(form.get(`${prefix}executable`)),
				args: parsed,
			},
		})
	}

	return result
}
