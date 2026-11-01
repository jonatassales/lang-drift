import type { Metadata } from 'next'
import { hasLocale } from 'next-intl'

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
