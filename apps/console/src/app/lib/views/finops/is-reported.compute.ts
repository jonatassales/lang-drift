export function isReported(value: number | null | undefined): value is number {
	return typeof value === 'number'
}
