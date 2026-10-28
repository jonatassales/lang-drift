import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
import { FinopsSection } from '@components/home/finops-section'
import { MicroruntimesSection } from '@components/home/microruntimes-section'
import { VisionSection } from '@components/home/vision-section'
import { PlatformHero } from '@components/platform-hero'
import { SkipLink } from '@components/skip-link'
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
			<SkipLink href="#top">{copy.skip}</SkipLink>
			<WebsiteHeader />
			<main id="top">
				<PlatformHero locale={locale} />
				<VisionSection locale={locale} />
				<MicroruntimesSection locale={locale} />
				<FinopsSection locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
