import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'

import { socialImages } from '@app/metadata/page-metadata'
import { earlyAccessSentMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface PageMetadataProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: PageMetadataProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { params } = props
	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const title = earlyAccessSentMessages[locale].title
	const images = await socialImages(parent)
	const openGraphImages = images.openGraph
	const twitterImages = images.twitter

	return {
		title: { absolute: title },
		robots: {
			index: false,
			follow: false,
		},
		openGraph: {
			title,
			...(openGraphImages ? { images: openGraphImages } : {}),
		},
		twitter: {
			card: 'summary_large_image',
			title,
			...(twitterImages ? { images: twitterImages } : {}),
		},
	}
}
