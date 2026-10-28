import { getLocale } from 'next-intl/server'

import { redirect } from '@i18n/navigation'

export default async function Page() {
	const locale = await getLocale()

	redirect({
		href: '/home',
		locale,
	})
}
