import { getLocale } from 'next-intl/server'

import { RegionRetry } from '@components/route-state/region-retry'
import type { SessionEnvelope } from '@data/read-envelopes.types'
import { readConsole } from '@data/server'
import type { User } from '@data/types'
import { redirect } from '@i18n/navigation'

import { AccountMenu } from './account-menu'
import { SignOutButton } from './sign-out-button'

export async function SessionControls() {
	let user: User | null
	try {
		const session = await readConsole<SessionEnvelope>('auth/session')

		user = session.user
	} catch {
		return <RegionRetry />
	}
	if (!user) {
		const locale = await getLocale()

		return redirect({
			href: '/login',
			locale,
		})
	}

	return (
		<AccountMenu email={user.email}>
			<SignOutButton />
		</AccountMenu>
	)
}
