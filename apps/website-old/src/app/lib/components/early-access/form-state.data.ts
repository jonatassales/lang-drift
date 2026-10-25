import type { EarlyAccessFormState } from './form.types'

export function createEarlyAccessFormState(): EarlyAccessFormState {
	return {
		email: '',
		hasInteractedWithEmail: false,
		fieldErrors: { email: null },
		status: 'idle',
	}
}
