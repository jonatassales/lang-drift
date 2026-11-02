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
