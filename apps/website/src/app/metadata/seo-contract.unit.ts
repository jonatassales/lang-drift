import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import type { Metadata } from 'next'
import { describe, expect, it } from 'vitest'

import routeState from '@i18n/messages/route-state.json'
import { routing } from '@i18n/routing'
import { generateMetadata as notFoundMetadata } from '@layout/not-found'

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

function robotsGroups(text: string) {
	return text
		.trim()
		.split(/\n\s*\n/)
		.map((block) => block.split('\n'))
		.filter((lines) => lines[0].startsWith('User-agent:'))
		.map((lines) => ({
			agents: lines
				.filter((line) => line.startsWith('User-agent:'))
				.map((line) => line.replace('User-agent:', '').trim()),
			rules: lines.filter((line) => !line.startsWith('User-agent:')),
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

describe('robots.txt', () => {
	const groups = robotsGroups(readWebsiteFile('src/app/robots.txt'))
	const trainingAgents = [
		'GPTBot',
		'ClaudeBot',
		'Google-Extended',
		'CCBot',
		'Applebot-Extended',
		'Bytespider',
	]

	it('lets every other crawler in, except on the API', () => {
		expect(groups.find((group) => group.agents.includes('*'))).toEqual({
			agents: ['*'],
			rules: ['Allow: /', 'Disallow: /api/'],
		})
	})

	it('keeps AI training crawlers out', () => {
		expect(groups.find((group) => group.agents.includes('GPTBot'))).toEqual({
			agents: trainingAgents,
			rules: ['Disallow: /'],
		})
	})

	it('names the sitemap on the public origin', () => {
		expect(readWebsiteFile('src/app/robots.txt')).toContain(
			`Sitemap: ${site.url}/sitemap.xml`,
		)
	})
})

describe('security.txt', () => {
	const fields = Object.fromEntries(
		readWebsiteFile('public/.well-known/security.txt')
			.trim()
			.split('\n')
			.map((line) => {
				const separator = line.indexOf(': ')

				return [line.slice(0, separator), line.slice(separator + 2)]
			}),
	)

	it('names the security contact, the languages and its canonical URL', () => {
		expect(fields).toMatchObject({
			Contact: 'mailto:security@langdrift.md',
			'Preferred-Languages': 'en, pt',
			Canonical: `${site.url}/.well-known/security.txt`,
		})
	})

	it('has not expired', () => {
		expect(new Date(fields.Expires).getTime()).toBeGreaterThan(Date.now())
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
