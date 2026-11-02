import type { Draft } from 'immer'

import type { EarlyAccessFormState, EmailErrorCode } from './form.types'

export function mapEmailBlur(
	draft: Draft<EarlyAccessFormState>,
	error: EmailErrorCode | null,
): void {
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = error
	if (error !== null) draft.status = 'idle'
}
