import type { Draft } from 'immer'

import type { EarlyAccessFormState, EmailErrorCode } from './form.types'

export function mapEmailChange(
	draft: Draft<EarlyAccessFormState>,
	email: string,
	error: EmailErrorCode | null,
): void {
	draft.email = email
	draft.fieldErrors.email = error
	draft.status = 'idle'
}

export function mapEmailBlur(
	draft: Draft<EarlyAccessFormState>,
	error: EmailErrorCode | null,
): void {
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = error
	if (error !== null) draft.status = 'idle'
}

export function mapValidationFailed(
	draft: Draft<EarlyAccessFormState>,
	error: EmailErrorCode,
): void {
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = error
	draft.status = 'idle'
}

export function mapValidationPassed(
	draft: Draft<EarlyAccessFormState>,
	email: string,
): void {
	draft.email = email
	draft.hasInteractedWithEmail = true
	draft.fieldErrors.email = null
	draft.status = 'verifying'
}
