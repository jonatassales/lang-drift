import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { PlatformFeatures } from '@components/platform-features'
import { PlatformHero } from '@components/platform-hero'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { platformMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface WebsitePageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: WebsitePageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('home', locale, await socialImages(parent))
}

export default async function WebsitePage(props: WebsitePageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = platformMessages[locale]

	return (
		<>
			<WebsiteStructuredData />
			<a
				className="fixed top-0 left-3 z-200 max-w-[calc(100vw-1.5rem)] -translate-y-full bg-foreground px-2.5 py-2 whitespace-normal text-card focus-visible:top-3 focus-visible:translate-y-0"
				href="#top"
			>
				{copy.skip}
			</a>
			<WebsiteHeader />
			<main id="top">
				<PlatformHero locale={locale} />
				<PlatformFeatures locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
