import { type NextRequest, NextResponse } from 'next/server'
import createMiddleware from 'next-intl/middleware'

import { cookieServiceFactory } from '@factories/cookie-service.factory'
import { localePath, routing, type WebsiteLocale } from '@i18n/routing'

const localeCookie = 'langdrift-website-locale'
const localeTtlSeconds = 15_552_000

const intlMiddleware = createMiddleware(routing)

const prefixes: Array<{ locale: WebsiteLocale; prefix: string }> = [
	{
		locale: 'pt-BR',
		prefix: '/pt-br',
	},
	{
		locale: 'zh-Hant',
		prefix: '/zh-hant',
	},
	{
		locale: 'ja',
		prefix: '/ja',
	},
	{
		locale: 'vi',
		prefix: '/vi',
	},
	{
		locale: 'en',
		prefix: '/en',
	},
]

export class LocaleService {
	handle(request: NextRequest): NextResponse {
		const segment = request.nextUrl.pathname.split('/')[1]
		const lowerSegment = segment.toLowerCase()
		const supported = prefixes.some(
			(entry) => entry.prefix === `/${lowerSegment}`,
		)
		const needsCanonicalRedirect = supported && segment !== lowerSegment
		if (needsCanonicalRedirect) {
			const canonical = new URL(request.url)
			canonical.pathname = canonical.pathname.replace(
				`/${segment}`,
				`/${lowerSegment}`,
			)

			return NextResponse.redirect(canonical, 308)
		}
		if (['zh-cn', 'zh-hans', 'zh-sg'].includes(lowerSegment)) {
			const fallback = new URL(request.url)
			fallback.pathname = fallback.pathname.replace(`/${segment}`, '/en')
			fallback.searchParams.set('availability', 'zh-CN')

			return NextResponse.redirect(fallback, 307)
		}
		const explicit = this.#explicitLocale(request.nextUrl.pathname)
		const reader = cookieServiceFactory(request)
		const saved = reader.get(localeCookie, (value): value is WebsiteLocale =>
			this.#isWebsiteLocale(value),
		)
		const browser = this.#matchAcceptLanguage(
			request.headers.get('accept-language'),
		)
		const unavailable = !explicit && !saved.found && browser === 'zh-CN'
		const locale = this.#resolveTargetLocale(explicit, saved, browser)
		const target = this.#localizedUrl(request, locale)
		if (unavailable) target.searchParams.set('availability', 'zh-CN')
		const response = explicit
			? this.#nextResponse(intlMiddleware(request))
			: NextResponse.redirect(target, 307)

		if (!explicit) {
			response.headers.set('Cache-Control', 'private, no-store')
			this.#mergeVary(response)
		}

		const shouldWrite =
			!this.#isPrefetch(request) && (!saved.found || saved.value !== locale)
		if (shouldWrite) {
			const writer = cookieServiceFactory(request, response)
			writer.set(localeCookie, locale, { ttlSeconds: localeTtlSeconds })
			response.headers.set('Cache-Control', 'private, no-store')
		}

		return response
	}

	#resolveTargetLocale(
		explicit: WebsiteLocale | undefined,
		saved: ReturnType<ReturnType<typeof cookieServiceFactory>['get']>,
		browser: WebsiteLocale | 'zh-CN' | undefined,
	): WebsiteLocale {
		if (explicit) return explicit
		const hasValidSavedLocale =
			saved.found && routing.locales.includes(saved.value as WebsiteLocale)
		if (hasValidSavedLocale) {
			return saved.value as WebsiteLocale
		}
		if (browser === 'zh-CN') return 'en'
		if (browser) return browser

		return routing.defaultLocale
	}

	#matchAcceptLanguage(
		header: string | null,
	): WebsiteLocale | 'zh-CN' | undefined {
		if (!header) return undefined

		const ranges = header
			.split(',')
			.map((part) => {
				const [tagPart, ...parameters] = part.trim().split(';')
				let quality = 1
				for (const parameter of parameters) {
					const [key, rawQuality] = parameter.trim().split('=')
					if (key !== 'q') continue

					const parsed = Number(rawQuality)
					quality = Number.isFinite(parsed) ? parsed : 0
				}

				return {
					tag: tagPart.trim().toLowerCase(),
					quality,
				}
			})
			.filter(
				(range) =>
					range.tag &&
					range.tag !== '*' &&
					range.quality > 0 &&
					range.quality <= 1,
			)
			.sort((left, right) => right.quality - left.quality)

		for (const range of ranges) {
			if (/^zh-(cn|hans|sg)(-|$)/.test(range.tag)) return 'zh-CN'
			const match = this.#matchLocaleTag(range.tag)
			if (match) return match
		}

		return undefined
	}

	#matchLocaleTag(tag: string): WebsiteLocale | undefined {
		const exact = routing.locales.find((locale) => locale.toLowerCase() === tag)
		if (exact) return exact
		if (tag === 'pt' || tag === 'pt-br' || tag.startsWith('pt-br-')) {
			return 'pt-BR'
		}
		if (
			tag === 'zh-hant' ||
			tag.startsWith('zh-hant-') ||
			tag === 'zh-tw' ||
			tag.startsWith('zh-tw-') ||
			tag === 'zh-hk' ||
			tag.startsWith('zh-hk-') ||
			tag === 'zh-mo' ||
			tag.startsWith('zh-mo-')
		) {
			return 'zh-Hant'
		}
		if (tag === 'en' || tag.startsWith('en-')) return 'en'
		if (tag === 'ja' || tag.startsWith('ja-')) return 'ja'
		if (tag === 'vi' || tag.startsWith('vi-')) return 'vi'

		return undefined
	}

	#isWebsiteLocale(value: unknown): value is WebsiteLocale {
		return routing.locales.some((locale) => locale === value)
	}

	#explicitLocale(pathname: string): WebsiteLocale | undefined {
		return prefixes.find(
			(entry) =>
				pathname === entry.prefix || pathname.startsWith(`${entry.prefix}/`),
		)?.locale
	}

	#localizedUrl(request: NextRequest, locale: WebsiteLocale): URL {
		const url = new URL(request.url)
		const pathname = request.nextUrl.pathname
		if (pathname === '/') {
			url.pathname = localePath[locale]

			return url
		}

		const prefix = localePath[locale].replace(/\/$/, '')
		const rest = pathname.startsWith('/') ? pathname : `/${pathname}`
		const joined = `${prefix}${rest}`
		url.pathname = joined.endsWith('/') ? joined : `${joined}/`

		return url
	}

	#isPrefetch(request: NextRequest): boolean {
		if (request.headers.has('next-router-prefetch')) return true
		if (request.headers.has('next-router-segment-prefetch')) return true

		return (
			this.#purposeHasPrefetch(request.headers.get('purpose')) ||
			this.#purposeHasPrefetch(request.headers.get('sec-purpose'))
		)
	}

	#purposeHasPrefetch(value: string | null): boolean {
		if (!value) return false

		return value
			.split(/[,;]/)
			.some((token) => token.trim().toLowerCase() === 'prefetch')
	}

	#mergeVary(response: NextResponse) {
		const current = response.headers.get('vary')
		const names = current
			? current
					.split(',')
					.map((part) => part.trim())
					.filter(Boolean)
			: []
		const present = new Set(names.map((name) => name.toLowerCase()))
		for (const name of ['Cookie', 'Accept-Language']) {
			if (!present.has(name.toLowerCase())) names.push(name)
		}
		response.headers.set('Vary', names.join(', '))
	}

	#nextResponse(response: Response): NextResponse {
		if (response instanceof NextResponse) return response

		throw new Error('Locale proxy expected a NextResponse')
	}
}
