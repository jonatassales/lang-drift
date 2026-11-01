'use client'

import { Button } from '@langdrift/web/ui/button'
import { useTranslations } from 'next-intl'
import { type FormEvent, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { mutation } from '@data/api'
import { Link, useRouter } from '@i18n/navigation'

import { AuthShell } from './auth-shell'

export function SignupView() {
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	const router = useRouter()
	const translate = useTranslations('ui')

	function handleSubmit(event: FormEvent<HTMLFormElement>) {
		event.preventDefault()
		const form = new FormData(event.currentTarget)
		setError(null)
		startTransition(async () => {
			try {
				await mutation('auth/signup', 'POST', {
					email: form.get('email'),
					password: form.get('password'),
				})

				router.replace('/home')
			} catch (failure) {
				setError(failure)
			}
		})
	}

	const submitLabel = pending ? translate('saving') : translate('signup')

	return (
		<AuthShell>
			<h1 className="sr-only">{translate('signup')}</h1>
			<form className="grid gap-5" method="post" onSubmit={handleSubmit}>
				<label className={labelClass}>
					{translate('email')}
					<input
						className={fieldClass}
						type="email"
						autoComplete="email"
						name="email"
						required
						maxLength={254}
					/>
				</label>
				<label className={labelClass}>
					{translate('password')}
					<input
						className={fieldClass}
						type="password"
						autoComplete="new-password"
						name="password"
						required
						minLength={12}
						maxLength={1024}
					/>
					<small className="font-normal text-muted-foreground">
						{translate('passwordHint')}
					</small>
				</label>
				{error ? <ErrorNotice error={error} /> : null}
				<Button type="submit" disabled={pending}>
					{submitLabel}
				</Button>
			</form>
			<Link
				href="/login"
				className="mt-5 block text-center text-xs text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
			>
				{translate('haveAccount')}
			</Link>
		</AuthShell>
	)
}
