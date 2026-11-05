// US dollars with cents, for figures from a fraction of a cent upwards.
export function workflowDollars(locale: string): Intl.NumberFormat {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 2,
		minimumFractionDigits: 2,
		style: 'currency',
	})
}
