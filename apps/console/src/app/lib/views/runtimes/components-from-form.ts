import type { Component } from '@data/types'

export function componentsFromForm(
	form: FormData,
	candidates: Component[],
): Component[] {
	const selected = new Set(form.getAll('selectedComponents').map(String))
	if (!selected.size) throw new Error('invalid_components')

	return candidates
		.filter((candidate) => selected.has(candidate.path))
		.map((candidate) => {
			const prefix = `component:${candidate.path}:`
			const parsed: unknown = JSON.parse(String(form.get(`${prefix}args`)))
			if (
				!Array.isArray(parsed) ||
				!parsed.every((argument) => typeof argument === 'string')
			)
				throw new Error('invalid_arguments')

			return {
				...candidate,
				image: String(form.get(`${prefix}image`)),
				port: Number(form.get(`${prefix}port`)),
				readinessPath: String(form.get(`${prefix}readinessPath`)),
				build: {
					executable: String(form.get(`${prefix}executable`)),
					args: parsed,
				},
			}
		})
}
