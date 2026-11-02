'use client'

import { produce } from 'immer'
import type { ChangeEvent, FocusEvent, FormEvent } from 'react'
import { useEffect, useId, useRef, useState } from 'react'

import type { EarlyAccessMessages } from '@i18n/messages/messages.types'
import { useRouter } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

import { handleEmailBlur } from './email-blur.handler'
import { handleEmailChange } from './email-change.handler'
import type { EarlyAccessFormState, StateUpdater } from './form.types'
import { createEarlyAccessFormState } from './form-state.data'
import { handleFormSubmit } from './form-submit.handler'
import type { TurnstileController } from './turnstile.controller'
import { createTurnstileController } from './turnstile.controller'

interface EarlyAccessFormProps {
	copy: EarlyAccessMessages
	helperId: string
	locale: WebsiteLocale
	source: 'landing' | 'pricing'
}

export function EarlyAccessForm(props: EarlyAccessFormProps) {
	const { copy, helperId, locale, source } = props

	const router = useRouter()
	const [formState, setFormState] = useState(createEarlyAccessFormState)
	const [hydrated, setHydrated] = useState(false)
	const inputRef = useRef<HTMLInputElement>(null)
	const submitting = useRef(false)
	const challenge = useRef<TurnstileController | null>(null)
	const mounted = useRef(false)
	const request = useRef<AbortController | null>(null)
	const transitioning = useRef(false)
	const challengeElement = useRef<HTMLDivElement>(null)
	const instanceId = useId()

	// Initialize Cloudflare Turnstile verification controller and clean up on unmount.
	useEffect(() => {
		mounted.current = true
		setHydrated(true)
		try {
			const siteKey = process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
			if (typeof siteKey !== 'string') throw new Error('unset')
			if (challengeElement.current) {
				challenge.current = createTurnstileController(
					challengeElement.current,
					siteKey,
				)
			}
		} catch {
			setFormState((state) => ({
				...state,
				status: 'challenge-error',
			}))
		}

		return () => {
			mounted.current = false
			request.current?.abort()
			challenge.current?.dispose()
			challenge.current = null
		}
	}, [])

	const emailError = formState.fieldErrors.email
	const update: StateUpdater<EarlyAccessFormState> = (recipe) => {
		if (mounted.current) setFormState(produce<EarlyAccessFormState>(recipe))
	}

	function handleEmailFieldBlur(event: FocusEvent<HTMLInputElement>) {
		handleEmailBlur(event.currentTarget.value, update)
	}

	function handleEmailFieldChange(event: ChangeEvent<HTMLInputElement>) {
		handleEmailChange(
			event.currentTarget.value,
			formState.hasInteractedWithEmail,
			update,
		)
	}

	function handleEmailFieldFocus() {
		void challenge.current?.prepare().catch(() => undefined)
	}

	function handleFormSubmission(event: FormEvent<HTMLFormElement>) {
		event.preventDefault()
		const controller = challenge.current
		if (submitting.current || !controller) return
		submitting.current = true
		request.current = new AbortController()

		function focusInput() {
			inputRef.current?.focus()
		}

		function navigateToSent() {
			if (!mounted.current) return
			transitioning.current = true
			try {
				router.push('/early-access-sent')
			} catch (error) {
				transitioning.current = false
				throw error
			}
		}

		void handleFormSubmit(
			inputRef.current?.value ?? '',
			update,
			focusInput,
			locale,
			source,
			controller,
			navigateToSent,
			request.current.signal,
		).finally(() => {
			if (!transitioning.current) submitting.current = false
			controller.reset()
		})
	}

	const busy =
		formState.status === 'verifying' ||
		formState.status === 'submitting' ||
		formState.status === 'success'
	const emailId = `${instanceId}-email`
	const errorId = `${instanceId}-error`
	const statusId = `${instanceId}-status`

	let statusText = ''
	if (formState.status === 'success') {
		statusText = copy.success
	} else if (
		formState.status === 'verifying' ||
		formState.status === 'submitting'
	) {
		statusText = copy.submitting
	} else if (formState.status === 'rate-limited') {
		statusText = copy.rateLimited.replace(
			'{seconds}',
			String(formState.retryAfter ?? 30),
		)
	} else if (formState.status === 'challenge-error') {
		statusText = copy.challengeFailed
	} else if (formState.status === 'error' && !emailError) {
		statusText = copy.unavailable
	}

	return (
		<form
			aria-busy={busy}
			className="min-w-0"
			noValidate
			onSubmit={handleFormSubmission}
		>
			<label
				className="mb-2.5 block text-sm font-medium text-foreground"
				htmlFor={emailId}
			>
				{copy.emailLabel}
			</label>
			<div className="flex flex-col gap-2.5 sm:flex-row sm:flex-wrap">
				<input
					aria-describedby={`${helperId}${emailError ? ` ${errorId}` : ''}`}
					aria-invalid={emailError !== null}
					autoComplete="email"
					className="min-h-12 w-full min-w-0 flex-initial rounded-[7px] border border-border bg-background px-3.5 py-3 text-base text-foreground placeholder:text-muted-foreground aria-[invalid=true]:border-status-danger disabled:cursor-wait disabled:opacity-65 sm:flex-[1_1_190px] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
					disabled={busy}
					id={emailId}
					inputMode="email"
					maxLength={254}
					onBlur={handleEmailFieldBlur}
					onChange={handleEmailFieldChange}
					onFocus={handleEmailFieldFocus}
					placeholder={copy.placeholder}
					ref={inputRef}
					required
					type="email"
					value={formState.email}
				/>
				<button
					className="inline-flex min-h-12 w-full cursor-pointer items-center justify-center gap-4.5 rounded-[7px] border border-primary bg-primary px-4.5 py-3 text-sm font-semibold leading-normal whitespace-nowrap text-primary-foreground hover:bg-primary-hover disabled:cursor-wait disabled:opacity-65 sm:w-auto focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
					disabled={!hydrated || busy}
					type="submit"
				>
					{busy ? copy.submitting : copy.submit}
					{busy ? (
						<span
							aria-hidden="true"
							className="inline-block size-4 animate-spin rounded-full border-2 border-current border-r-transparent motion-reduce:animate-none"
						/>
					) : (
						<span aria-hidden="true" className="shrink-0">
							↗
						</span>
					)}
				</button>
			</div>
			<div ref={challengeElement} />
			<output
				aria-atomic="true"
				aria-live="polite"
				className="block min-h-12 text-[13px] leading-[1.6] [overflow-wrap:anywhere] [&_p]:mt-2.5 [&_p]:mb-0"
			>
				{emailError ? (
					<span className="text-status-danger" id={errorId}>
						{copy.errors[emailError]}
					</span>
				) : null}
				<span className="text-foreground" id={statusId}>
					{statusText}
				</span>
			</output>
		</form>
	)
}
