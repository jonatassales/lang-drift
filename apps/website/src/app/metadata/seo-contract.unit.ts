import { describe, expect, it } from 'vitest'

import { canonicalOf } from '@app/metadata/testdata/canonical-of.compute'
import { languagesOf } from '@app/metadata/testdata/languages-of.compute'
import { readWebsiteFile } from '@app/metadata/testdata/read-website-file.mock'
import { robotsGroups } from '@app/metadata/testdata/robots-groups.compute'
import { sitemapEntries } from '@app/metadata/testdata/sitemap-entries.compute'
import routeState from '@i18n/messages/route-state.json'
import { routing } from '@i18n/routing'
import { generateMetadata as notFoundMetadata } from '@layout/not-found'

import locales from './locales.json'
import { pageMetadata } from './page-metadata.map'
import pages from './pages.json'
import site from './site.json'

const pageIds = Object.keys(pages) as Array<keyof typeof pages>

const localizedPages = pageIds.flatMap((page) =>
	routing.locales.map((locale) => ({
		page,
		locale,
		metadata: pageMetadata(page, locale),
	})),
)

const sitemap = sitemapEntries(readWebsiteFile('src/app/sitemap.xml'))

describe('website/metadata/seo-contract', () => {
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

	it.each(routing.locales)('has its title in %s', async (locale) => {
		const metadata = await notFoundMetadata({
			params: Promise.resolve({ locale }),
		})

		expect(metadata.title).toBe(routeState[locale].notFound)
	})
})
