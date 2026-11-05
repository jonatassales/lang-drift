import { describe, expect, it } from 'vitest'

import type { WebsiteLocale } from '@i18n/routing'

import { assertCatalogShape } from './assert-catalog-shape.compute'

const catalog = {
	en: {
		title: 1,
	},
	'pt-BR': {
		title: 1,
	},
	'zh-Hant': {
		title: 1,
	},
	ja: {
		title: 1,
	},
	vi: {
		title: 1,
	},
}

describe('website/template/i18n/loaders/assert-catalog-shape.compute', () => {
	it('returns a catalog whose locales share the English shape', () => {
		expect(assertCatalogShape('sample', catalog)).toBe(catalog)
	})

	it('refuses a locale that does not match the English shape', () => {
		const broken: Record<WebsiteLocale, unknown> = {
			...catalog,
			ja: {
				heading: 1,
			},
		}

		expect(() => assertCatalogShape('sample', broken)).toThrow(
			'sample locale ja does not match the catalog shape',
		)
	})
})
