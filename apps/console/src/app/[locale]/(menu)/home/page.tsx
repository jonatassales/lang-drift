import { getTranslations } from 'next-intl/server'
import { Suspense } from 'react'

import { RuntimeListSkeleton } from '@components/route-state/runtime-list-skeleton'
import { readConsole } from '@data/server'
import { AccountHome } from '@views/home/account-home'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('home') }
}

export default async function Page() {
	return (
		<Suspense fallback={<RuntimeListSkeleton />}>
			<AccountHome initial={readConsole('home')} />
		</Suspense>
	)
}
