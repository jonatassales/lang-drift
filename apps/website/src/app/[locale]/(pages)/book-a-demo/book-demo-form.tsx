'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useId, useRef, useState } from 'react'
import { useForm } from 'react-hook-form'

import { isBookDemoErrorCode } from '@api/book-a-demo/is-book-demo-error-code.validation'
import type { BookDemoMessages } from '@i18n/messages/messages.types'
import type { WebsiteLocale } from '@i18n/routing'
import { Toast } from 'web/ui/toast'

import { DemoNotice } from './demo-notice'
import type {
	BookDemoData,
	BookDemoField,
	BookDemoInput,
	BookDemoOutcome,
} from './form.types'
import { bookDemoSchema } from './form.validation'
import { bookDemoFormState } from './form-state.data'
import { positionIds } from './positions'

interface BookDemoFormProps {
	copy: BookDemoMessages
	locale: WebsiteLocale
}

const label = 'text-sm font-medium text-foreground'
const fieldControl =
	'min-h-11 w-full rounded-lg border border-border bg-background px-3.5 text-foreground aria-[invalid=true]:border-status-danger'
const inputControl = `${fieldControl} focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary`
const selectControl = `${fieldControl} text-sm`
const errorText = 'm-0 text-xs text-status-danger'

// The demo request: validated as the visitor leaves each field and again on
// submit, posted with the locale to the edge, which forwards it to the
// demo-requests service, and answered with the return-soon notice while that
// service is unreachable.
export function BookDemoForm(props: BookDemoFormProps) {
	const { copy, locale } = props

	const [outcome, setOutcome] = useState<BookDemoOutcome>('idle')
	const [errorMessage, setErrorMessage] = useState<string | null>(null)
	const noticeRef = useRef<HTMLHeadingElement>(null)
	const fieldId = useId()

	const { formState, handleSubmit, register } = useForm<
		BookDemoInput,
		unknown,
		BookDemoData
	>({
		defaultValues: bookDemoFormState,
		mode: 'onTouched',
		resolver: zodResolver(bookDemoSchema),
	})

	// Focus the notice heading once the request is sent or found unavailable,
	// so a screen reader announces the outcome that replaced the form.
	useEffect(() => {
		if (outcome !== 'idle') noticeRef.current?.focus()
	}, [outcome])

	function errorOf(field: BookDemoField) {
		const code = formState.errors[field]?.message
		if (!isBookDemoErrorCode(code)) return undefined

		return copy.errors[code]
	}

	function handleDismissError() {
		setErrorMessage(null)
	}

	async function send(data: BookDemoData) {
		setErrorMessage(null)

		try {
			const response = await fetch('/api/book-a-demo/', {
				body: JSON.stringify({
					...data,
					locale,
				}),
				headers: { 'content-type': 'application/json' },
				method: 'POST',
			})

			if (response.ok) {
				setOutcome('sent')

				return
			}

			if (response.status === 503) {
				setOutcome('unavailable')

				return
			}

			setErrorMessage(copy.form.failed)
		} catch {
			setErrorMessage(copy.form.failed)
		}
	}

	if (outcome !== 'idle') {
		return (
			<div>
				<DemoNotice copy={copy.form[outcome]} headingRef={noticeRef} />
			</div>
		)
	}

	const emailError = errorOf('email')
	const companyError = errorOf('company')
	const positionError = errorOf('position')

	return (
		<form
			aria-busy={formState.isSubmitting}
			className="grid gap-5"
			noValidate
			onSubmit={handleSubmit(send)}
		>
			<div className="grid gap-2">
				<label className={label} htmlFor={`${fieldId}-email`}>
					{copy.form.email}
				</label>
				<input
					aria-describedby={emailError ? `${fieldId}-email-error` : undefined}
					aria-invalid={Boolean(emailError)}
					autoComplete="email"
					className={inputControl}
					id={`${fieldId}-email`}
					inputMode="email"
					placeholder={copy.form.emailPlaceholder}
					type="email"
					{...register('email')}
				/>
				{emailError ? (
					<p className={errorText} id={`${fieldId}-email-error`}>
						{emailError}
					</p>
				) : null}
			</div>
			<div className="grid gap-2">
				<label className={label} htmlFor={`${fieldId}-company`}>
					{copy.form.company}
				</label>
				<input
					aria-describedby={
						companyError ? `${fieldId}-company-error` : undefined
					}
					aria-invalid={Boolean(companyError)}
					autoComplete="organization"
					className={inputControl}
					id={`${fieldId}-company`}
					placeholder={copy.form.companyPlaceholder}
					type="text"
					{...register('company')}
				/>
				{companyError ? (
					<p className={errorText} id={`${fieldId}-company-error`}>
						{companyError}
					</p>
				) : null}
			</div>
			<div className="grid gap-2">
				<label className={label} htmlFor={`${fieldId}-position`}>
					{copy.form.position}
				</label>
				<select
					aria-describedby={
						positionError ? `${fieldId}-position-error` : undefined
					}
					aria-invalid={Boolean(positionError)}
					autoComplete="organization-title"
					className={selectControl}
					id={`${fieldId}-position`}
					{...register('position')}
				>
					<option value="">{copy.form.positionPlaceholder}</option>
					{positionIds.map((position) => (
						<option key={position} value={position}>
							{copy.positions[position]}
						</option>
					))}
				</select>
				{positionError ? (
					<p className={errorText} id={`${fieldId}-position-error`}>
						{positionError}
					</p>
				) : null}
			</div>
			<div className="flex flex-wrap items-center gap-4">
				<button
					className="smooth-primary-button whitespace-nowrap focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
					disabled={formState.isSubmitting}
					type="submit"
				>
					{formState.isSubmitting ? copy.form.submitting : copy.form.submit}
				</button>
			</div>
			{errorMessage ? (
				<Toast onClose={handleDismissError}>{errorMessage}</Toast>
			) : null}
		</form>
	)
}
