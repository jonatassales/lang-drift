import { mapEmailChange } from './email-change.map'
import type { EarlyAccessFormState, StateUpdater } from './form.types'
import { getEmailError } from './form.validation'

export function handleEmailChange(
	email: string,
	hasInteractedWithEmail: boolean,
	update: StateUpdater<EarlyAccessFormState>,
): void {
	const error = hasInteractedWithEmail ? getEmailError(email) : null
	update((draft) => {
		mapEmailChange(draft, email, error)
	})
}
