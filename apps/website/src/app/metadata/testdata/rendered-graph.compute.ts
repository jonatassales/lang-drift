import { pageStructuredData } from '@app/metadata/page-structured-data.map'
import type pages from '@app/metadata/pages.json'
import { serializeJsonLd } from '@app/metadata/serialize-json-ld.compute'
import type { routing } from '@i18n/routing'

import type { JsonLdGraph } from './json-ld-graph.types'

export function renderedGraph(
	page: keyof typeof pages,
	locale: (typeof routing.locales)[number],
): JsonLdGraph {
	return JSON.parse(serializeJsonLd(pageStructuredData(page, locale)))
}
