export function currency(
	value: number,
	currency: string,
	locale = 'en-US',
): string {
	return new Intl.NumberFormat(locale, {
		style: 'currency',
		currency,
	}).format(value)
}
