// The ticket widgets' figures in the viewer's locale: token counts from thousands,
// compact as the website shows them.
export function compactTokens(locale: string, thousands: number) {
	return new Intl.NumberFormat(locale, {
		notation: 'compact',
	}).format(thousands * 1000)
}
