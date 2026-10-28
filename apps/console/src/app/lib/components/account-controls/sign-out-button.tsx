'use client'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { mutation } from '@data/api'
import { useRouter } from '@i18n/navigation'
import { SignOut } from '@ui/icons'

export function SignOutButton() {
	const t = useTranslations('ui')
	const router = useRouter()
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	return (
		<div className="relative">
			<button
				disabled={pending}
				aria-label={t('logout')}
				title={t('logout')}
				type="button"
				className="flex h-8 w-8 items-center justify-center rounded-md hover:bg-muted"
				onClick={() =>
					startTransition(async () => {
						setError(null)
						try {
							await mutation('auth/logout', 'POST')

							router.replace('/login')
							router.refresh()
						} catch (failure) {
							setError(failure)
						}
					})
				}
			>
				<SignOut aria-hidden size={16} />
			</button>
			{error ? (
				<div className="absolute right-0 top-full w-64">
					<ErrorNotice error={error} />
				</div>
			) : null}
		</div>
	)
}
