import type { WebsiteLocale } from '@i18n/routing'

import organization from './organization.json'
import { type PageId, pageUrl } from './page-metadata'
import pages from './pages.json'
import site from './site.json'

const organizationId = `${site.url}/#organization`
const websiteId = `${site.url}/#website`

const organizationNode = {
	'@id': organizationId,
	'@type': 'Organization',
	name: organization.name,
	url: organization.url,
	logo: organization.logo,
	sameAs: organization.sameAs,
}

const websiteNode = {
	'@id': websiteId,
	'@type': 'WebSite',
	name: site.name,
	publisher: { '@id': organizationId },
	url: `${site.url}/`,
}

export function pageStructuredData(page: PageId, locale: WebsiteLocale) {
	const url = pageUrl(page, locale)
	const copy = pages[page][locale]
	const webPage = {
		'@id': `${url}#webpage`,
		'@type': 'WebPage',
		url,
		name: copy.title,
		description: copy.description,
		inLanguage: locale,
		isPartOf: { '@id': websiteId },
	}
	const siteNodes = page === 'home' ? [organizationNode, websiteNode] : []

	return {
		'@context': 'https://schema.org',
		'@graph': [...siteNodes, webPage],
	}
}
