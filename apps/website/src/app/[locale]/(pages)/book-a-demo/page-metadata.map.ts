import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'

import { pageMetadata } from '@app/metadata/page-metadata.map'
import { socialImages } from '@app/metadata/social-images.adapter'
import { routing } from '@i18n/routing'

interface BookDemoPageMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: BookDemoPageMetadataProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const images = await socialImages(parent)

	return pageMetadata('book-a-demo', locale, images)
}
