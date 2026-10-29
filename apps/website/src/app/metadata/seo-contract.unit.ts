import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import type { Metadata, ResolvingMetadata } from 'next'
import { describe, expect, it } from 'vitest'

import { generateMetadata as earlyAccessSentMetadata } from '@app/[locale]/(pages)/early-access-sent/page'
import { generateMetadata as notFoundMetadata } from '@app/[locale]/not-found'
import routeState from '@i18n/messages/route-state.json'
import { routing } from '@i18n/routing'

import locales from './locales.json'
import { pageMetadata } from './page-metadata'
import pages from './pages.json'
import site from './site.json'

interface SitemapEntry {
	loc: string
	alternates: Record<string, string>
}

const pageIds = Object.keys(pages) as Array<keyof typeof pages>

const localizedPages = pageIds.flatMap((page) =>
	routing.locales.map((locale) => ({
		page,
		locale,
		metadata: pageMetadata(page, locale),
	})),
)

const sitemap = sitemapEntries(readWebsiteFile('src/app/sitemap.xml'))

function readWebsiteFile(path: string) {
	return readFileSync(join(process.cwd(), path), 'utf8')
}

function sitemapEntries(xml: string): SitemapEntry[] {
	return Array.from(xml.matchAll(/<url>([\s\S]*?)<\/url>/g), ([, body]) => ({
		loc: body.match(/<loc>(.*?)<\/loc>/)?.[1] ?? '',
		alternates: Object.fromEntries(
			Array.from(
				body.matchAll(/hreflang="([^"]+)" href="([^"]+)"/g),
				([, language, href]) => [language, href],
			),
		),
	}))
}

function canonicalOf(metadata: Metadata) {
	return String(metadata.alternates?.canonical)
}

function languagesOf(metadata: Metadata) {
	return Object.fromEntries(
		Object.entries(metadata.alternates?.languages ?? {}).map(
			([language, href]) => [language, String(href)],
		),
	)
}

describe('Every indexable page in every locale', () => {
	it.each(routing.locales)(
		'has a unique title and description in %s',
		(locale) => {
			const copies = pageIds.map((page) => pages[page][locale])

			expect(new Set(copies.map((copy) => copy.title)).size).toBe(
				pageIds.length,
			)
			expect(new Set(copies.map((copy) => copy.description)).size).toBe(
				pageIds.length,
			)
		},
	)

	it.each(localizedPages)(
		'$page in $locale carries its title and description',
		({ page, locale, metadata }) => {
			expect(metadata.title).toEqual({ absolute: pages[page][locale].title })
			expect(metadata.description).toBe(pages[page][locale].description)
		},
	)

	it.each(localizedPages)(
		'$page in $locale has a self canonical on the public origin',
		({ metadata }) => {
			const canonical = new URL(canonicalOf(metadata))

			expect(canonical.origin).toBe(site.url)
			expect(canonical.pathname.endsWith('/')).toBe(true)
			expect(canonical.search).toBe('')
		},
	)

	it.each(localizedPages)(
		'$page in $locale names every locale and its English URL as x-default',
		({ page, metadata }) => {
			const languages = Object.fromEntries(
				routing.locales.map((language) => [
					language,
					canonicalOf(pageMetadata(page, language)),
				]),
			)

			expect(languagesOf(metadata)).toEqual({
				...languages,
				'x-default': languages.en,
			})
		},
	)

	it.each(localizedPages)(
		'$page in $locale shares its canonical through Open Graph and Twitter',
		({ page, locale, metadata }) => {
			const otherLocales = routing.locales
				.filter((language) => language !== locale)
				.map((language) => locales[language].openGraphLocale)

			expect(metadata.openGraph).toMatchObject({
				url: canonicalOf(metadata),
				title: pages[page][locale].title,
				description: pages[page][locale].description,
				locale: locales[locale].openGraphLocale,
				alternateLocale: otherLocales,
			})
			expect(metadata.twitter).toMatchObject({
				card: 'summary_large_image',
				title: pages[page][locale].title,
			})
		},
	)

	it.each(localizedPages)('$page in $locale is indexable', ({ metadata }) => {
		expect(metadata.robots).toMatchObject({
			index: true,
			follow: true,
		})
	})
})

describe('The sitemap', () => {
	it('lists exactly the indexable pages', () => {
		const canonicals = localizedPages.map(({ metadata }) =>
			canonicalOf(metadata),
		)

		expect(sitemap.map((entry) => entry.loc).sort()).toEqual(canonicals.sort())
	})

	it.each(localizedPages)(
		'gives $page in $locale the alternates of its HTML',
		({ metadata }) => {
			const entry = sitemap.find((item) => item.loc === canonicalOf(metadata))

			expect(entry?.alternates).toEqual(languagesOf(metadata))
		},
	)
})

describe('The early-access confirmation page', () => {
	const parent = Promise.resolve({
		openGraph: null,
		twitter: null,
	}) as unknown as ResolvingMetadata

	it.each(routing.locales)('stays out of search in %s', async (locale) => {
		const metadata = await earlyAccessSentMetadata(
			{ params: Promise.resolve({ locale }) },
			parent,
		)

		expect(metadata.robots).toMatchObject({ index: false })
		expect(
			sitemap.some((entry) => entry.loc.includes('early-access-sent')),
		).toBe(false)
	})
})

describe('The localized 404', () => {
	it.each(routing.locales)('has its title in %s', async (locale) => {
		const metadata = await notFoundMetadata({
			params: Promise.resolve({ locale }),
		})

		expect(metadata.title).toBe(routeState[locale].notFound)
	})
})
