import { describe, expect, it } from 'vitest'

import { renderedGraph } from '@app/metadata/testdata/rendered-graph.compute'
import { routing } from '@i18n/routing'

import organization from './organization.json'
import { pageMetadata } from './page-metadata'
import { serializeJsonLd } from './page-structured-data.map'
import pages from './pages.json'
import site from './site.json'

const pageIds = Object.keys(pages) as Array<keyof typeof pages>

const localizedPages = pageIds.flatMap((page) =>
	routing.locales.map((locale) => ({
		page,
		locale,
	})),
)

describe('website/metadata/page-structured-data.map', () => {
	it.each(localizedPages)(
		'describes $page in $locale as one WebPage of the website',
		({ page, locale }) => {
			const graph = renderedGraph(page, locale)
			const canonical = String(pageMetadata(page, locale).alternates?.canonical)

			expect(graph['@context']).toBe('https://schema.org')
			expect(
				graph['@graph'].filter((node) => node['@type'] === 'WebPage'),
			).toEqual([
				{
					'@id': `${canonical}#webpage`,
					'@type': 'WebPage',
					url: canonical,
					name: pages[page][locale].title,
					description: pages[page][locale].description,
					inLanguage: locale,
					isPartOf: { '@id': `${site.url}/#website` },
				},
			])
		},
	)

	it.each(routing.locales)(
		'keeps the Organization and WebSite nodes on the home page in %s',
		(locale) => {
			const graph = renderedGraph('home', locale)

			expect(graph['@graph']).toEqual(
				expect.arrayContaining([
					{
						'@id': `${site.url}/#organization`,
						'@type': 'Organization',
						name: organization.name,
						url: organization.url,
						logo: organization.logo,
						sameAs: organization.sameAs,
					},
					{
						'@id': `${site.url}/#website`,
						'@type': 'WebSite',
						name: site.name,
						publisher: { '@id': `${site.url}/#organization` },
						url: `${site.url}/`,
					},
				]),
			)
		},
	)

	it('escapes < so no value can close the script element', () => {
		const value = { name: '</script><script>alert(1)</script>' }
		const json = serializeJsonLd(value)

		expect(json).not.toContain('<')
		expect(JSON.parse(json)).toEqual(value)
	})
})
