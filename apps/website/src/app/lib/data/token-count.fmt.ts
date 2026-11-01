// Thousands of tokens as a compact count: 3.5M, 250K.
export function tokenCount(locale: string, thousands: number): string {
	return new Intl.NumberFormat(locale, {
		maximumFractionDigits: 1,
		notation: 'compact',
	}).format(thousands * 1000)
}
