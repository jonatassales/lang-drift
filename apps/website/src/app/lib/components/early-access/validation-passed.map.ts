import type { Draft } from 'immer'

import type { EarlyAccessFormState } from './form.types'

export function mapValidationPassed(
	draft: Draft<EarlyAccessFormState>,
	email: string,
): void {
	draft.email = email
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = null
	draft.status = 'verifying'
}
