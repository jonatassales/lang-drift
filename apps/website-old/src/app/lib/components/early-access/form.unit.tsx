// @vitest-environment jsdom
import {
	act,
	cleanup,
	fireEvent,
	render,
	screen,
	waitFor,
} from '@testing-library/react'
import { StrictMode } from 'react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./form.submission', () => ({ submitEarlyAccess: vi.fn() }))
const push = vi.hoisted(() => vi.fn())
vi.mock('@i18n/navigation', () => ({ useRouter: () => ({ push }) }))

import { earlyAccessMessages } from '@i18n/messages/index'

import { EarlyAccessForm } from './form'
import { submitEarlyAccess } from './form.submission'
import type { TurnstileOptions } from './turnstile.types'

let callbacks: TurnstileOptions[]
beforeEach(() => {
	vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', '1x00000000000000000000BB')
	vi.clearAllMocks()
	callbacks = []
	window.turnstile = {
		render: vi.fn((_element, options) => {
			callbacks.push(options)
			return `widget-${callbacks.length}`
		}),
		execute: vi.fn(),
		isExpired: vi.fn(() => false),
		reset: vi.fn(),
		remove: vi.fn(),
	}
	vi.mocked(submitEarlyAccess).mockResolvedValue({ ok: true })
})
afterEach(() => {
	cleanup()
	delete window.turnstile
	vi.unstubAllEnvs()
})
function mount() {
	return render(
		<StrictMode>
			<EarlyAccessForm
				copy={earlyAccessMessages.en}
				helperId="helper"
				locale="en"
				source="landing"
			/>
		</StrictMode>,
	)
}
function submit(email = 'a@example.com') {
	fireEvent.change(screen.getByLabelText('Email address'), {
		target: { value: email },
	})
	fireEvent.submit(
		screen.getByRole('button').closest('form') as HTMLFormElement,
	)
}
async function executed() {
	await waitFor(() => expect(window.turnstile?.execute).toHaveBeenCalled())
}

describe('invisible Early Access form', () => {
	it('has no challenge chrome and prepares on focus without executing', async () => {
		const { container } = mount()
		expect(
			container.querySelector(
				'fieldset, legend, [name="cf-turnstile-response"]',
			),
		).toBeNull()
		expect(window.turnstile?.render).not.toHaveBeenCalled()
		fireEvent.focus(screen.getByLabelText('Email address'))
		await waitFor(() => expect(callbacks).toHaveLength(1))
		expect(callbacks[0]).toMatchObject({
			execution: 'execute',
			'response-field': false,
			retry: 'never',
			'refresh-expired': 'never',
		})
		expect(window.turnstile?.execute).not.toHaveBeenCalled()
		expect(submitEarlyAccess).not.toHaveBeenCalled()
	})

	it('validates and focuses an invalid email without executing', () => {
		mount()
		submit('invalid')
		expect(
			screen.getByText(earlyAccessMessages.en.errors.emailInvalid),
		).toBeTruthy()
		expect(document.activeElement).toBe(screen.getByLabelText('Email address'))
		expect(window.turnstile?.execute).not.toHaveBeenCalled()
		expect(submitEarlyAccess).not.toHaveBeenCalled()
	})

	it('waits for a callback and suppresses duplicates through server completion', async () => {
		let finish!: (value: { ok: true }) => void
		vi.mocked(submitEarlyAccess).mockReturnValue(
			new Promise((resolve) => {
				finish = resolve
			}),
		)
		mount()
		submit()
		submit()
		await executed()
		expect(window.turnstile?.execute).toHaveBeenCalledTimes(1)
		expect(submitEarlyAccess).not.toHaveBeenCalled()
		expect((screen.getByRole('button') as HTMLButtonElement).disabled).toBe(
			true,
		)
		await act(async () => {
			callbacks[0]?.callback('first-token')
			callbacks[0]?.callback('duplicate-token')
		})
		submit()
		expect(submitEarlyAccess).toHaveBeenCalledExactlyOnceWith(
			{
				email: 'a@example.com',
				locale: 'en',
				source: 'landing',
				turnstileToken: 'first-token',
			},
			expect.any(AbortSignal),
		)
		await act(async () => {
			finish({ ok: true })
		})
		expect(screen.getByText(earlyAccessMessages.en.success)).toBeTruthy()
		expect(
			(screen.getByLabelText('Email address') as HTMLInputElement).value,
		).toBe('a@example.com')
		expect(window.turnstile?.reset).toHaveBeenCalledWith('widget-1')
		expect(window.turnstile?.remove).toHaveBeenCalledWith('widget-1')
	})

	it.each(['error', 'expired', 'unsupported'] as const)(
		'recovers from %s and ignores an old callback on retry',
		async (failure) => {
			mount()
			submit()
			await executed()
			const old = callbacks[0]
			if (!old) throw new Error('Missing prepared widget')
			await act(async () => {
				if (failure === 'error')
					old['error-callback']('private-provider-detail')
				if (failure === 'expired') old['expired-callback']()
				if (failure === 'unsupported') old['unsupported-callback']()
			})
			expect(
				screen.getByText(earlyAccessMessages.en.challengeFailed),
			).toBeTruthy()
			expect(submitEarlyAccess).not.toHaveBeenCalled()
			submit()
			await waitFor(() => expect(callbacks).toHaveLength(2))
			await act(async () => {
				old.callback('stale-token')
				callbacks[1]?.callback('fresh-token')
			})
			expect(submitEarlyAccess).toHaveBeenCalledExactlyOnceWith(
				expect.objectContaining({ turnstileToken: 'fresh-token' }),
				expect.any(AbortSignal),
			)
		},
	)

	it('discards a token expired between callback and dispatch', async () => {
		mount()
		submit()
		await executed()
		await act(async () => {
			callbacks[0]?.callback('expired-token')
			callbacks[0]?.['expired-callback']()
		})
		expect(submitEarlyAccess).not.toHaveBeenCalled()
		expect(
			screen.getByText(earlyAccessMessages.en.challengeFailed),
		).toBeTruthy()
	})

	it.each(['CHALLENGE_FAILED', 'UNAVAILABLE'] as const)(
		're-executes after server %s without reusing a token',
		async (code) => {
			vi.mocked(submitEarlyAccess).mockResolvedValueOnce({
				ok: false,
				code,
			})
			mount()
			submit()
			await executed()
			await act(async () => {
				callbacks[0]?.callback('spent-token')
			})
			expect(
				(screen.getByLabelText('Email address') as HTMLInputElement).value,
			).toBe('a@example.com')
			submit()
			await waitFor(() => expect(callbacks).toHaveLength(2))
			await act(async () => {
				callbacks[1]?.callback('retry-token')
			})
			expect(window.turnstile?.execute).toHaveBeenCalledTimes(2)
			expect(submitEarlyAccess).toHaveBeenNthCalledWith(
				2,
				expect.objectContaining({ turnstileToken: 'retry-token' }),
				expect.any(AbortSignal),
			)
		},
	)

	it('removes its widget and prevents submission after unmount', async () => {
		const view = mount()
		submit()
		await executed()
		view.unmount()
		await act(async () => {
			callbacks[0]?.callback('late-token')
		})
		expect(submitEarlyAccess).not.toHaveBeenCalled()
		expect(window.turnstile?.remove).toHaveBeenCalledWith('widget-1')
	})
})

it('announces rate limiting and preserves the address', async () => {
	vi.mocked(submitEarlyAccess).mockResolvedValue({
		ok: false,
		code: 'RATE_LIMITED',
		retryAfter: 60,
	})
	mount()
	submit()
	await executed()
	await act(async () => {
		callbacks[0]?.callback('token')
	})
	expect(screen.getByRole('status').textContent).toContain('60')
	expect(
		(screen.getByLabelText('Email address') as HTMLInputElement).value,
	).toBe('a@example.com')
	expect((screen.getByRole('button') as HTMLButtonElement).disabled).toBe(false)
})
it('keeps controls locked through success navigation', async () => {
	mount()
	submit()
	await executed()
	expect(screen.getByRole('status').textContent).toContain(
		earlyAccessMessages.en.submitting,
	)
	await act(async () => {
		callbacks[0]?.callback('token')
	})
	submit()
	expect(submitEarlyAccess).toHaveBeenCalledTimes(1)
	expect(push).toHaveBeenCalledExactlyOnceWith('/early-access-sent')
	expect((screen.getByRole('button') as HTMLButtonElement).disabled).toBe(true)
})
it('aborts an in-flight request and never navigates after unmount', async () => {
	let complete!: (result: { ok: true }) => void
	vi.mocked(submitEarlyAccess).mockReturnValue(
		new Promise((resolve) => {
			complete = resolve
		}),
	)
	const view = mount()
	submit()
	await executed()
	await act(async () => {
		callbacks[0]?.callback('token')
	})
	const signal = vi.mocked(submitEarlyAccess).mock.calls[0]?.[1]
	view.unmount()
	expect(signal?.aborted).toBe(true)
	await act(async () => {
		complete({ ok: true })
	})
	expect(push).not.toHaveBeenCalled()
})
