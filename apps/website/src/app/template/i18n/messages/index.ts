import type { ExecutiveReviewMessages } from '@data/executive-review-data'
import type { FinopsMessages } from '@data/finops-content'
import type { MicroruntimesMessages } from '@data/microruntimes-content'
import type { VisionLoopMessages } from '@data/vision-loop-content'
import type { WebsiteLocale } from '@i18n/routing'

import conversionJson from './conversion.json'
import demoJson from './demo.json'
import earlyAccessJson from './early-access.json'
import earlyAccessSentJson from './early-access-sent.json'
import reviewEn from './executive-review/en.json'
import reviewJa from './executive-review/ja.json'
import reviewPt from './executive-review/pt-br.json'
import reviewVi from './executive-review/vi.json'
import reviewZh from './executive-review/zh-hant.json'
import finopsJson from './finops.json'
import footerJson from './footer.json'
import headerJson from './header.json'
import homeJson from './home.json'
import integrationsJson from './integrations.json'
import type {
	ConversionCopy,
	DemoMessages,
	EarlyAccessMessages,
	EarlyAccessSentMessages,
	FooterMessages,
	HeaderMessages,
	HomeMessages,
	IntegrationsMessages,
	PlatformMessages,
	PrivacyMessages,
} from './messages.types'
import microruntimesJson from './microruntimes.json'
import platformJson from './platform.json'
import privacyJson from './privacy.json'
import routeStateJson from './route-state.json'
import visionEn from './vision-loop/en.json'
import visionJa from './vision-loop/ja.json'
import visionPt from './vision-loop/pt-br.json'
import visionVi from './vision-loop/vi.json'
import visionZh from './vision-loop/zh-hant.json'

const locales = ['en', 'pt-BR', 'zh-Hant', 'ja', 'vi'] as const

function shape(value: unknown): string {
	if (Array.isArray(value))
		return `[${value.map((item) => shape(item)).join(',')}]`
	if (value !== null && typeof value === 'object') {
		return `{${Object.keys(value)
			.map((key) => `${key}:${shape(Reflect.get(value, key))}`)
			.join(',')}}`
	}

	return typeof value
}

function assertCatalogShape(
	name: string,
	catalog: Record<WebsiteLocale, unknown>,
) {
	const base = shape(catalog.en)
	for (const locale of locales) {
		if (shape(catalog[locale]) !== base) {
			throw new Error(
				`${name} locale ${locale} does not match the catalog shape`,
			)
		}
	}
}

function triple(
	value: readonly string[],
	label: string,
): [string, string, string] {
	if (value.length !== 3) {
		throw new Error(`${label} must contain exactly three strings`)
	}
	const first = value[0]
	const second = value[1]
	const third = value[2]
	if (
		typeof first !== 'string' ||
		typeof second !== 'string' ||
		typeof third !== 'string'
	) {
		throw new Error(`${label} must contain exactly three strings`)
	}

	return [first, second, third]
}

function loadConversion(locale: (typeof conversionJson)['en']): ConversionCopy {
	return {
		...locale,
		roi: {
			...locale.roi,
			brief: {
				...locale.roi.brief,
				evidence: triple(locale.roi.brief.evidence, 'roi.brief.evidence'),
			},
		},
		plans: {
			...locale.plans,
			cards: {
				trial: {
					...locale.plans.cards.trial,
					items: triple(
						locale.plans.cards.trial.items,
						'plans.cards.trial.items',
					),
				},
				plus: {
					...locale.plans.cards.plus,
					items: triple(
						locale.plans.cards.plus.items,
						'plans.cards.plus.items',
					),
				},
				pro: {
					...locale.plans.cards.pro,
					items: triple(locale.plans.cards.pro.items, 'plans.cards.pro.items'),
				},
			},
		},
	}
}

assertCatalogShape('home', homeJson)
assertCatalogShape('header', headerJson)
assertCatalogShape('demo', demoJson)
assertCatalogShape('privacy', privacyJson)
assertCatalogShape('conversion', conversionJson)
assertCatalogShape('early-access', earlyAccessJson)
assertCatalogShape('early-access-sent', earlyAccessSentJson)
assertCatalogShape('integrations', integrationsJson)
assertCatalogShape('footer', footerJson)
assertCatalogShape('platform', platformJson)
assertCatalogShape('route-state', routeStateJson)
assertCatalogShape('finops', finopsJson)
assertCatalogShape('microruntimes', microruntimesJson)

export const homeMessages: Record<WebsiteLocale, HomeMessages> = homeJson
export const headerMessages: Record<WebsiteLocale, HeaderMessages> = headerJson
export const demoMessages: Record<WebsiteLocale, DemoMessages> = demoJson
export const privacyMessages: Record<WebsiteLocale, PrivacyMessages> =
	privacyJson
export const earlyAccessMessages: Record<WebsiteLocale, EarlyAccessMessages> =
	earlyAccessJson
export const earlyAccessSentMessages: Record<
	WebsiteLocale,
	EarlyAccessSentMessages
> = earlyAccessSentJson
export const integrationsMessages: Record<WebsiteLocale, IntegrationsMessages> =
	integrationsJson
export const footerMessages: Record<WebsiteLocale, FooterMessages> = footerJson
export const platformMessages: Record<WebsiteLocale, PlatformMessages> =
	platformJson
export const finopsMessages: Record<WebsiteLocale, FinopsMessages> = finopsJson
export const microruntimesMessages: Record<
	WebsiteLocale,
	MicroruntimesMessages
> = microruntimesJson

export const conversionMessages: Record<WebsiteLocale, ConversionCopy> = {
	en: loadConversion(conversionJson.en),
	'pt-BR': loadConversion(conversionJson['pt-BR']),
	'zh-Hant': loadConversion(conversionJson['zh-Hant']),
	ja: loadConversion(conversionJson.ja),
	vi: loadConversion(conversionJson.vi),
}

const visionLoopCatalog = {
	en: visionEn,
	'pt-BR': visionPt,
	'zh-Hant': visionZh,
	ja: visionJa,
	vi: visionVi,
}
assertCatalogShape('vision-loop', visionLoopCatalog)
export const visionLoopMessages: Record<WebsiteLocale, VisionLoopMessages> =
	visionLoopCatalog

const executiveReviewCatalog = {
	en: reviewEn,
	'pt-BR': reviewPt,
	'zh-Hant': reviewZh,
	ja: reviewJa,
	vi: reviewVi,
}
assertCatalogShape('executive-review', executiveReviewCatalog)
export const executiveReviewMessages: Record<
	WebsiteLocale,
	ExecutiveReviewMessages
> = executiveReviewCatalog

export function getWebsiteMessages(locale: WebsiteLocale) {
	return {
		home: homeMessages[locale],
		header: headerMessages[locale],
		demo: demoMessages[locale],
		routeState: routeStateJson[locale],
	}
}
