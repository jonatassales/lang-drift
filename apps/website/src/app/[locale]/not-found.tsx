import type { Metadata } from 'next'
import { hasLocale } from 'next-intl'
import { getLocale, getTranslations } from 'next-intl/server'

import { RouteState } from '@components/route-state'
import { getWebsiteMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface NotFoundMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: NotFoundMetadataProps,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) return {}

	return { title: getWebsiteMessages(locale).routeState.notFound }
}

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
			imageSrc="/route-state/not-found.svg"
			imageAlt={translate('notFoundImage')}
			locale={locale}
			title={translate('notFound')}
		/>
	)
}
