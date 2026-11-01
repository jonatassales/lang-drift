import type { Metadata, ResolvingMetadata } from 'next'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { pageMetadata, socialImages } from '@app/metadata/page-metadata'
// import { EarlyAccessSection } from '@components/early-access/section'
import { PlansSection } from '@components/plans-section'
import { SkipLink } from '@components/skip-link'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { WebsiteStructuredData } from '@components/website-structured-data'
import { conversionMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

interface PricingPageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: PricingPageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()

	return pageMetadata('pricing', locale, await socialImages(parent))
}

export default async function PricingPage(props: PricingPageProps) {
	const { params } = props

	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)

	return (
		<>
			<WebsiteStructuredData locale={locale} page="pricing" />
			<SkipLink href="#plans-main">
				{conversionMessages[locale].plans.skip}
			</SkipLink>
			<WebsiteHeader />
			<main className="min-w-0" id="plans-main">
				<PlansSection locale={locale} />
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
