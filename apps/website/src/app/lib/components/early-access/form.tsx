'use client'

import { produce } from 'immer'
import { useEffect, useId, useRef, useState } from 'react'

import type { EarlyAccessMessages } from '@i18n/messages/messages.types'
import { useRouter } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

import {
	handleEmailBlur,
	handleEmailChange,
	handleFormSubmit,
} from './form.handlers'
import type { EarlyAccessFormState, StateUpdater } from './form.types'
import { createEarlyAccessFormState } from './form-state.data'
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
	const busy =
		formState.status === 'verifying' ||
		formState.status === 'submitting' ||
		formState.status === 'success'
	const instanceId = useId()
	const emailId = `${instanceId}-email`
	const errorId = `${instanceId}-error`
	const statusId = `${instanceId}-status`
	const emailError = formState.fieldErrors.email
	const update: StateUpdater<EarlyAccessFormState> = (recipe) => {
		if (mounted.current) setFormState(produce<EarlyAccessFormState>(recipe))
	}

	return (
		<form
			className="early-access-form"
			noValidate
			aria-busy={busy}
			onSubmit={(event) => {
				event.preventDefault()
				const controller = challenge.current
				if (submitting.current || !controller) return
				submitting.current = true
				request.current = new AbortController()
				void handleFormSubmit(
					inputRef.current?.value ?? '',
					update,
					() => {
						inputRef.current?.focus()
					},
					locale,
					source,
					controller,
					() => {
						if (!mounted.current) return
						transitioning.current = true
						try {
							router.push('/early-access-sent')
						} catch (error) {
							transitioning.current = false
							throw error
						}
					},
					request.current.signal,
				).finally(() => {
					if (!transitioning.current) submitting.current = false
					controller.reset()
				})
			}}
		>
			<label htmlFor={emailId}>{copy.emailLabel}</label>
			<div className="early-access-controls">
				<input
					className="focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
					ref={inputRef}
					id={emailId}
					type="email"
					autoComplete="email"
					inputMode="email"
					required
					maxLength={254}
					disabled={busy}
					value={formState.email}
					placeholder={copy.placeholder}
					aria-describedby={`${helperId}${emailError ? ` ${errorId}` : ''}`}
					aria-invalid={emailError !== null}
					onChange={(event) => {
						handleEmailChange(
							event.currentTarget.value,
							formState.hasInteractedWithEmail,
							update,
						)
					}}
					onFocus={() => {
						void challenge.current?.prepare().catch(() => {
							// A submission retries preparation and reports a recoverable error.
						})
					}}
					onBlur={(event) => {
						handleEmailBlur(event.currentTarget.value, update)
					}}
				/>
				<button
					className="whitespace-nowrap focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
					type="submit"
					disabled={!hydrated || busy}
				>
					{busy ? copy.submitting : copy.submit}
					{busy ? (
						<span
							aria-hidden="true"
							className="inline-block size-4 animate-spin rounded-full border-2 border-current border-r-transparent motion-reduce:animate-none"
						/>
					) : (
						<span aria-hidden="true">↗</span>
					)}
				</button>
			</div>
			<div ref={challengeElement} />
			<output
				className="early-access-feedback"
				aria-live="polite"
				aria-atomic="true"
			>
				{emailError ? (
					<span id={errorId} className="early-access-field-error">
						{copy.errors[emailError]}
					</span>
				) : null}
				<span id={statusId} className="early-access-valid-status">
					{formState.status === 'success'
						? copy.success
						: formState.status === 'verifying' ||
								formState.status === 'submitting'
							? copy.submitting
							: formState.status === 'rate-limited'
								? copy.rateLimited.replace(
										'{seconds}',
										String(formState.retryAfter ?? 30),
									)
								: formState.status === 'challenge-error'
									? copy.challengeFailed
									: formState.status === 'error' && !emailError
										? copy.unavailable
										: ''}
				</span>
			</output>
		</form>
	)
}
