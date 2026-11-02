import type { Metadata } from 'next'

import { routing, type WebsiteLocale } from '@i18n/routing'

import {
	type InheritedSocialImages,
	inheritedImages,
} from './inherited-images.compute'
import locales from './locales.json'
import type { PageId } from './page-path.compute'
import { pageUrl } from './page-url.compute'
import pages from './pages.json'
import site from './site.json'

export { pageUrl } from './page-url.compute'
export { socialImages } from './social-images.compute'
export type { InheritedSocialImages, PageId }

const indexRobots: Metadata['robots'] = {
	index: true,
	follow: true,
	googleBot: {
		index: true,
		follow: true,
		'max-image-preview': 'large',
		'max-snippet': -1,
		'max-video-preview': -1,
	},
}

export function pageMetadata(
	page: PageId,
	locale: WebsiteLocale,
	images?: InheritedSocialImages,
): Metadata {
	const copy = pages[page][locale]
	const canonical = pageUrl(page, locale)
	const languages = Object.fromEntries(
		routing.locales.map((language) => [language, pageUrl(page, language)]),
	)
	const alternateLocale = routing.locales
		.filter((language) => language !== locale)
		.map((language) => locales[language].openGraphLocale)
	const openGraphImages = inheritedImages(images?.openGraph)
	const twitterImages = inheritedImages(images?.twitter)

	return {
		title: { absolute: copy.title },
		description: copy.description,
		alternates: {
			canonical,
			languages: {
				...languages,
				'x-default': pageUrl(page, routing.defaultLocale),
			},
		},
		openGraph: {
			type: 'website',
			siteName: site.name,
			title: copy.title,
			description: copy.description,
			url: canonical,
			locale: locales[locale].openGraphLocale,
			alternateLocale,
			...(openGraphImages ? { images: openGraphImages } : {}),
		},
		twitter: {
			card: 'summary_large_image',
			title: copy.title,
			description: copy.description,
			...(twitterImages ? { images: twitterImages } : {}),
		},
		robots: indexRobots,
	}
}
