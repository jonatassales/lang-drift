import { mapEmailBlur } from './email-blur.map'
import type { EarlyAccessFormState, StateUpdater } from './form.types'
import { getEmailError } from './form.validation'

export function handleEmailBlur(
	email: string,
	update: StateUpdater<EarlyAccessFormState>,
): void {
	const error = getEmailError(email)
	update((draft) => {
		mapEmailBlur(draft, error)
	})
}
