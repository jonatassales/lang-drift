import { Suspense } from 'react'

import { RuntimeListSkeleton } from '@components/route-state/runtime-list-skeleton'
import { readConsole } from '@data/server'
import type { HomeOverview } from '@data/types'
import { redirect } from '@i18n/navigation'
import { AccountHome } from '@views/home/account-home'

export { generateMetadata } from './page.metadata'

interface PageProps {
	params: Promise<{ locale: string }>
}

export default async function Page(props: PageProps) {
	const { params } = props

	const { locale } = await params

	const initial = readConsole<HomeOverview>('home')

	const home = await initial

	if (home.products.length > 0 && home.running.length === 0) {
		redirect({
			href: `/products/${home.products[0].id}`,
			locale,
		})
	}

	return (
		<Suspense fallback={<RuntimeListSkeleton />}>
			<AccountHome initial={initial} />
		</Suspense>
	)
}
