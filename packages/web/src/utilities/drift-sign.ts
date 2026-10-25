export function driftSign(value: number): 'negative' | 'positive' | 'zero' {
	return value < 0 ? 'negative' : value > 0 ? 'positive' : 'zero'
}

export function driftColor(value: number): string {
	switch (driftSign(value)) {
		case 'negative':
			return 'var(--lg-drift-negative)'
		case 'positive':
			return 'var(--lg-drift-positive)'
		case 'zero':
			return 'var(--lg-drift-zero)'
	}
}
