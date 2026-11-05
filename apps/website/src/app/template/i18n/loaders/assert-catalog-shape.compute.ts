import { routing, type WebsiteLocale } from '@i18n/routing'

import { shape } from './shape.compute'

// Every locale of a catalog carries the English catalog's shape: the same
// keys, nested the same way, with values of the same kind. A loader passes
// each catalog through it as it assembles the catalog, so a broken locale
// fails as soon as the loader loads.
export function assertCatalogShape<Catalog>(
	name: string,
	catalog: Record<WebsiteLocale, Catalog>,
): Record<WebsiteLocale, Catalog> {
	const base = shape(catalog.en)

	for (const locale of routing.locales) {
		if (shape(catalog[locale]) !== base) {
			throw new Error(
				`${name} locale ${locale} does not match the catalog shape`,
			)
		}
	}

	return catalog
}
