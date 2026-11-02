import { hasLocale } from 'next-intl'
import { getLocale, getTranslations } from 'next-intl/server'

import { RouteState } from '@components/route-state'
import { routing } from '@i18n/routing'

export { generateMetadata } from './not-found.metadata'

export default async function NotFound() {
	const locale = await getLocale()

	if (!hasLocale(routing.locales, locale)) {
		throw new Error(`Unsupported request locale: ${locale}`)
	}

	const translate = await getTranslations('routeState')

	return (
		<RouteState
			description={translate('notFoundDescription')}
			homeLabel={translate('home')}
			imageAlt={translate('notFoundImage')}
			imageSrc="/route-state/not-found.svg"
			locale={locale}
			title={translate('notFound')}
		/>
	)
}
