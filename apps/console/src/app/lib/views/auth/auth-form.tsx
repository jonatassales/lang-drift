'use client'
import Image from 'next/image'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { StubBanner } from '@components/source-notice/stub-banner'
import { mutation } from '@data/api'
import { LanguageSettings } from '@i18n/language-settings'
import { Link, useRouter } from '@i18n/navigation'
import { Button } from '@ui/button'
import { Card } from '@ui/card'

interface AuthFormProps {
	signup?: boolean
	source: string
}
export function AuthForm(props: AuthFormProps) {
	const { signup = false, source } = props

	const t = useTranslations('ui')
	const router = useRouter()
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	return (
		<div className="mx-auto flex min-h-screen max-w-md flex-col justify-center gap-6 px-5 py-10">
			<div className="text-center">
				<Image
					src="/brand/exphere-logo.svg"
					alt="Exphere"
					width={165}
					height={38}
					priority
					className="mx-auto h-auto w-[165px] dark:invert"
				/>
				<h1 className="mt-4 text-2xl font-medium">
					{t(signup ? 'signup' : 'login')}
				</h1>
				<p className="mt-2 text-sm text-muted-foreground">
					{t(signup ? 'signupIntro' : 'loginIntro')}
				</p>
			</div>
			<StubBanner source={source} />
			<Card className="p-6">
				<form
					className="grid gap-5"
					onSubmit={(event) => {
						event.preventDefault()
						const form = new FormData(event.currentTarget)
						setError(null)
						startTransition(async () => {
							try {
								await mutation(`auth/${signup ? 'signup' : 'login'}`, 'POST', {
									email: form.get('email'),
									password: form.get('password'),
								})

								router.replace('/home')
							} catch (failure) {
								setError(failure)
							}
						})
					}}
				>
					<label className={labelClass}>
						{t('email')}
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
						{t('password')}
						<input
							className={fieldClass}
							type="password"
							autoComplete={signup ? 'new-password' : 'current-password'}
							name="password"
							required
							minLength={signup ? 12 : 1}
							maxLength={1024}
						/>
						{signup ? (
							<small className="font-normal text-muted-foreground">
								{t('passwordHint')}
							</small>
						) : null}
					</label>
					{error ? <ErrorNotice error={error} /> : null}
					<Button type="submit" tone="primary" disabled={pending}>
						{t(pending ? 'saving' : signup ? 'signup' : 'login')}
					</Button>
				</form>
				<Link
					href={signup ? '/login' : '/signup'}
					className="mt-5 block text-center text-xs text-primary"
				>
					{t(signup ? 'haveAccount' : 'needAccount')}
				</Link>
			</Card>
			<div className="flex justify-center">
				<LanguageSettings />
			</div>
		</div>
	)
}
