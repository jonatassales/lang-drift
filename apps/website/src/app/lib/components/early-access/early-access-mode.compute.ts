export type EarlyAccessMode = 'production' | 'development' | 'test'

export function earlyAccessMode(): EarlyAccessMode {
	if (process.env.NODE_ENV === 'development') return 'development'
	if (process.env.EARLY_ACCESS_MODE === 'test') return 'test'

	return 'production'
}
