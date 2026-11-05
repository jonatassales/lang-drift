// Amounts on the FinOps page are simulated US dollars, whole only.
export function dollars(locale: string) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 0,
		style: 'currency',
	})
}
