// The ticket widgets' figures in the viewer's locale: dollars from cents.
export function dollars(locale: string, cents: number) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		style: 'currency',
	}).format(cents / 100)
}
