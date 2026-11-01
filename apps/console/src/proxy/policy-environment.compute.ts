import type { PolicyEnvironment } from './content-security-policy'

export function policyEnvironment(): PolicyEnvironment {
	if (process.env.NODE_ENV === 'development') return 'development'
	if (process.env.VERCEL_ENV === 'preview') return 'preview'

	return 'production'
}
