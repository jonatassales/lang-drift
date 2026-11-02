import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { routing } from '@i18n/routing'

interface PrivacyPageMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: PrivacyPageMetadataProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const images = await socialImages(parent)

	return pageMetadata('privacy', locale, images)
}
