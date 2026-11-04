'use client'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { mutation } from '@data/api'
import { useRouter } from '@i18n/navigation'
import { SignOut } from 'web/ui/icons/client'

export function SignOutButton() {
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	const router = useRouter()
	const translate = useTranslations('ui')

	function handleSignOut() {
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

	return (
		<div className="relative">
			<button
				disabled={pending}
				type="button"
				className="flex min-h-9 w-full cursor-pointer items-center gap-2.5 rounded-md px-2.5 text-left text-[12px] text-foreground hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
				onClick={handleSignOut}
			>
				<SignOut aria-hidden size={15} />
				{translate('logout')}
			</button>
			{error ? (
				<div className="absolute right-0 top-full w-64">
					<ErrorNotice error={error} />
				</div>
			) : null}
		</div>
	)
}
