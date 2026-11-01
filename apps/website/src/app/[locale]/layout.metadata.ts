import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'

import site from '@app/metadata/site.json'
import { routing } from '@i18n/routing'

interface LayoutMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: LayoutMetadataProps,
): Promise<Metadata> {
	const { params } = props
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return {
		metadataBase: new URL(site.url),
		applicationName: site.name,
	}
}
