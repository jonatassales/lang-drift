import type { Metadata, ResolvingMetadata } from 'next'

import { localePath, routing, type WebsiteLocale } from '@i18n/routing'

import locales from './locales.json'
import pages from './pages.json'
import site from './site.json'

type PageId = keyof typeof pages

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

// Pages whose URL is not their id, so the metadata JSON stores no route path.
const nestedPagePath: Partial<Record<PageId, string>> = {
	finops: 'platform/finops/',
	microruntimes: 'platform/microruntimes/',
}

function pagePath(page: PageId, locale: WebsiteLocale) {
	if (page === 'home') return localePath[locale]

	return `${localePath[locale]}${nestedPagePath[page] ?? `${page}/`}`
}

function absolute(pathname: string) {
	return new URL(pathname, site.url).toString()
}

interface InheritedSocialImages {
	openGraph?: NonNullable<Awaited<ResolvingMetadata>['openGraph']>['images']
	twitter?: NonNullable<Awaited<ResolvingMetadata>['twitter']>['images']
}

function inheritedImages(
	images: InheritedSocialImages['openGraph'] | InheritedSocialImages['twitter'],
) {
	if (!images) return
	if (Array.isArray(images) && images.length === 0) return

	return images
}

export async function socialImages(parent: ResolvingMetadata) {
	const resolved = await parent

	return {
		openGraph: resolved.openGraph?.images,
		twitter: resolved.twitter?.images,
	}
}

export function pageMetadata(
	page: PageId,
	locale: WebsiteLocale,
	images?: InheritedSocialImages,
): Metadata {
	const copy = pages[page][locale]
	const canonical = absolute(pagePath(page, locale))
	const languages = Object.fromEntries(
		routing.locales.map((language) => [
			language,
			absolute(pagePath(page, language)),
		]),
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
				'x-default': absolute('/'),
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
