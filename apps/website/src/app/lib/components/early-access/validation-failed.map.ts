import type { Draft } from 'immer'

import type { EarlyAccessFormState, EmailErrorCode } from './form.types'

export function mapValidationFailed(
	draft: Draft<EarlyAccessFormState>,
	error: EmailErrorCode,
): void {
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = error
	draft.status = 'idle'
}
