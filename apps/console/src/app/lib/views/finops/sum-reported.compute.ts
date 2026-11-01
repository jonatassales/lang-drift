import { isReported } from './is-reported.compute.ts'

export function sumReported(
	values: (number | null | undefined)[],
): number | null {
	const reported = values.filter(isReported)
	if (reported.length === 0) return null

	return reported.reduce((total, value) => total + value, 0)
}
