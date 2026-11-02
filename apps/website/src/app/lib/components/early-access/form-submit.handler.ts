import type { WebsiteLocale } from '@i18n/routing'

import { submitEarlyAccess } from './form.submission'
import type { EarlyAccessFormState, StateUpdater } from './form.types'
import { earlyAccessFormSchema, getEmailError } from './form.validation'
import type { SubmissionChallenge } from './turnstile.types'
import { mapValidationFailed } from './validation-failed.map'
import { mapValidationPassed } from './validation-passed.map'

export async function handleFormSubmit(
	email: string,
	update: StateUpdater<EarlyAccessFormState>,
	focusEmail: () => void,
	locale: WebsiteLocale,
	source: 'landing' | 'pricing',
	challenge: SubmissionChallenge,
	onSuccess?: () => void,
	signal?: AbortSignal,
): Promise<void> {
	const result = earlyAccessFormSchema.safeParse({ email })
	if (!result.success) {
		const error = getEmailError(email) ?? 'emailInvalid'
		update((draft) => {
			mapValidationFailed(draft, error)
		})
		focusEmail()

		return
	}

	update((draft) => {
		mapValidationPassed(draft, result.data.email)
	})

	let turnstileToken: string
	try {
		await challenge.execute()

		signal?.throwIfAborted()
		turnstileToken = challenge.takeToken()
	} catch {
		update((draft) => {
			draft.status = 'challenge-error'
		})

		return
	}

	update((draft) => {
		draft.status = 'submitting'
	})

	try {
		const response = await submitEarlyAccess(
			{
				...result.data,
				locale,
				source,
				turnstileToken,
			},
			signal,
		)

		signal?.throwIfAborted()
		update((draft) => {
			draft.status = response.ok
				? 'success'
				: response.code === 'CHALLENGE_FAILED'
					? 'challenge-error'
					: response.code === 'RATE_LIMITED'
						? 'rate-limited'
						: 'error'
			if (!response.ok && response.code === 'RATE_LIMITED')
				draft.retryAfter = response.retryAfter
			if (!response.ok && response.code === 'VALIDATION_ERROR') {
				const error = response.fieldErrors?.email
				draft.fieldErrors.email =
					error === 'emailRequired' ||
					error === 'emailTooLong' ||
					error === 'emailInvalid'
						? error
						: null
			}
		})
		if (response.ok) onSuccess?.()
	} catch {
		update((draft) => {
			draft.status = 'error'
		})
	}
}
