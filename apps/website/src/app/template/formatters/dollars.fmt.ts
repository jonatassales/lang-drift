// Amounts on the LLMOps and FinOps pages are US dollars, whole only.
export function dollars(locale: string) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 0,
		style: 'currency',
	})
}
