import { produce } from 'immer'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import type {
	EarlyAccessActionResult,
	EarlyAccessFormState,
	StateUpdater,
} from './form.types'
import { createEarlyAccessFormState } from './form-state.data'

vi.mock('./form.submission', () => ({ submitEarlyAccess: vi.fn() }))

import { handleEmailBlur } from './email-blur.handler'
import { handleEmailChange } from './email-change.handler'
import { submitEarlyAccess } from './form.submission'
import { handleFormSubmit } from './form-submit.handler'

const challenge = {
	execute: vi.fn(async () => {}),
	takeToken: vi.fn(() => 'token'),
}
let state: EarlyAccessFormState
const update: StateUpdater<EarlyAccessFormState> = (recipe) => {
	state = produce(state, recipe)
}

beforeEach(() => {
	state = createEarlyAccessFormState()
	vi.clearAllMocks()
})

describe('form transitions', () => {
	it('validates on blur and on changes after interaction', () => {
		handleEmailChange('bad', false, update)
		expect(state.fieldErrors.email).toBeNull()
		handleEmailBlur('bad', update)
		expect(state.fieldErrors.email).toBe('emailInvalid')
		handleEmailChange('a@example.com', true, update)
		expect(state.fieldErrors.email).toBeNull()
		expect(state.status).toBe('idle')
	})

	it('focuses invalid email and does not call the server', async () => {
		const focus = vi.fn()
		await handleFormSubmit('bad', update, focus, 'en', 'landing', challenge)
		expect(focus).toHaveBeenCalledOnce()
		expect(state.fieldErrors.email).toBe('emailInvalid')
		expect(submitEarlyAccess).not.toHaveBeenCalled()
	})

	it.each([
		[{ ok: true }, 'success', 'a@example.com'],
		[
			{
				ok: false,
				code: 'CHALLENGE_FAILED',
			},
			'challenge-error',
			'a@example.com',
		],
		[
			{
				ok: false,
				code: 'UNAVAILABLE',
			},
			'error',
			'a@example.com',
		],
		[
			{
				ok: false,
				code: 'VALIDATION_ERROR',
				fieldErrors: { email: 'emailTooLong' },
			},
			'error',
			'a@example.com',
		],
		[
			{
				ok: false,
				code: 'VALIDATION_ERROR',
				fieldErrors: { email: 'bogus' },
			},
			'error',
			'a@example.com',
		],
	] as const)(
		'maps server result %j to %s',
		async (response, status, email) => {
			vi.mocked(submitEarlyAccess).mockResolvedValue(
				response as EarlyAccessActionResult,
			)
			await handleFormSubmit(
				' a@example.com ',
				update,
				vi.fn(),
				'en',
				'pricing',
				challenge,
			)
			expect(state.status).toBe(status)
			expect(state.email).toBe(email)
			expect(submitEarlyAccess).toHaveBeenCalledWith(
				{
					email: 'a@example.com',
					locale: 'en',
					source: 'pricing',
					turnstileToken: 'token',
				},
				undefined,
			)
			if (!response.ok && response.code === 'VALIDATION_ERROR') {
				expect(state.fieldErrors.email).toBe(
					response.fieldErrors.email === 'emailTooLong' ? 'emailTooLong' : null,
				)
			}
		},
	)

	it('recovers from an unexpected action rejection', async () => {
		vi.mocked(submitEarlyAccess).mockRejectedValue(new Error('offline'))
		await handleFormSubmit(
			'a@example.com',
			update,
			vi.fn(),
			'en',
			'landing',
			challenge,
		)
		expect(state.status).toBe('error')
	})

	it('navigates only after the server accepts the registration', async () => {
		const onSuccess = vi.fn()
		vi.mocked(submitEarlyAccess).mockResolvedValueOnce({ ok: true })
		await handleFormSubmit(
			'a@example.com',
			update,
			vi.fn(),
			'en',
			'landing',
			challenge,
			onSuccess,
		)
		expect(onSuccess).toHaveBeenCalledOnce()

		vi.mocked(submitEarlyAccess).mockResolvedValueOnce({
			ok: false,
			code: 'UNAVAILABLE',
		})
		await handleFormSubmit(
			'a@example.com',
			update,
			vi.fn(),
			'en',
			'landing',
			challenge,
			onSuccess,
		)
		expect(onSuccess).toHaveBeenCalledOnce()
	})
})
