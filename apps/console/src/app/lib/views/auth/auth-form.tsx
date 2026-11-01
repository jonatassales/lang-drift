'use client'
import { Brand } from '@langdrift/web/ui/brand'
import { Button } from '@langdrift/web/ui/button'
import { ElevatedFrame } from '@langdrift/web/ui/elevated-frame'
import { HeroShader } from '@langdrift/web/ui/hero-shader'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ThemeControl } from '@components/account-controls/theme-control'
import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { mutation } from '@data/api'
import { LanguageSettings } from '@i18n/language-settings'
import { Link, useRouter } from '@i18n/navigation'

interface AuthFormProps {
	signup?: boolean
}
// The console's way in, on the website's ground: the shader in its
// muted tone across the screen, the language at the top left and the
// theme at the top right as in the website's header, the LangDrift brand alone
// above the form, and the form in the elevated frame.
export function AuthForm(props: AuthFormProps) {
	const { signup = false } = props

	const t = useTranslations('ui')
	const router = useRouter()
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	return (
		<section className="relative isolate min-h-screen overflow-hidden bg-background">
			<HeroShader tone="muted" />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="absolute inset-x-0 top-0 z-10 flex items-center justify-between px-5 py-4 sm:px-7">
				<div className="relative h-8 w-14">
					<div className="absolute top-1/2 left-0 -translate-y-1/2">
						<LanguageSettings compact />
					</div>
				</div>
				<ThemeControl />
			</div>
			<div className="relative mx-auto flex min-h-screen max-w-md flex-col justify-center gap-6 px-5 py-10">
				<div className="flex justify-center">
					<Brand className="text-foreground" />
				</div>
				<h1 className="sr-only">{t(signup ? 'signup' : 'login')}</h1>
				<ElevatedFrame className="p-6">
					<form
						className="grid gap-5"
						method="post"
						onSubmit={(event) => {
							event.preventDefault()
							const form = new FormData(event.currentTarget)
							setError(null)
							startTransition(async () => {
								try {
									await mutation(
										`auth/${signup ? 'signup' : 'login'}`,
										'POST',
										{
											email: form.get('email'),
											password: form.get('password'),
										},
									)

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
						<Button type="submit" disabled={pending}>
							{t(pending ? 'saving' : signup ? 'signup' : 'login')}
						</Button>
					</form>
					<Link
						href={signup ? '/login' : '/signup'}
						className="mt-5 block text-center text-xs text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
					>
						{t(signup ? 'haveAccount' : 'needAccount')}
					</Link>
				</ElevatedFrame>
			</div>
		</section>
	)
}
