// The ticket widgets' figures in the viewer's locale: dollars from cents and
// token counts from thousands, compact as the website shows them.
export function dollars(locale: string, cents: number) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		style: 'currency',
	}).format(cents / 100)
}

export function compactTokens(locale: string, thousands: number) {
	return new Intl.NumberFormat(locale, {
		notation: 'compact',
	}).format(thousands * 1000)
}
