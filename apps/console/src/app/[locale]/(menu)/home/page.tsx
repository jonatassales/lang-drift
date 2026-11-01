import { Suspense } from 'react'

import { RuntimeListSkeleton } from '@components/route-state/runtime-list-skeleton'
import { readConsole } from '@data/server'
import { AccountHome } from '@views/home/account-home'

export { generateMetadata } from './page.metadata'

export default async function Page() {
	return (
		<Suspense fallback={<RuntimeListSkeleton />}>
			<AccountHome initial={readConsole('home')} />
		</Suspense>
	)
}
